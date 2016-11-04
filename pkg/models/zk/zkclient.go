package zkclient

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/anywhy/medis/pkg/utils/errors"
	"github.com/anywhy/medis/pkg/utils/log"
	"github.com/samuel/go-zookeeper/zk"
	"path/filepath"
)

var ErrClientClosed = errors.New("use a closed zookeeper client")

var DefaultLogger = func(format string, v ...interface{}) {
	log.Infof("zookeeper - " + fmt.Sprintf(format, v...))
}

type zkLogger struct {
	logger func(format string, v ...interface{})
}

func (z *zkLogger) Printf(format string, v ...interface{}) {
	if nil != z && z.logger != nil {
		z.logger(format, v...)
	}
}

type Client struct {
	sync.Mutex
	conn *zk.Conn

	zkAddr  string
	timeout time.Duration

	logger *zkLogger
	dialAt time.Time
	closed bool
}

func New(addrs string, timeout time.Duration) (*Client, error) {
	return NewZkClient(addrs, timeout, DefaultLogger)
}

func NewZkClient(zkAddr string, timeout time.Duration, logger func(foramt string, v ...interface{})) (*Client, error) {
	if timeout <= 0 {
		timeout = time.Second * 30
	}

	client := &Client{
		zkAddr:  zkAddr,
		timeout: timeout,
	}

	if err := client.reset(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) reset() error {
	c.dialAt = time.Now()
	conn, events, err := zk.Connect(strings.Split(c.zkAddr, ","), c.timeout)
	if err != nil {
		return errors.Trace(err)
	}
	if c.conn != nil {
		c.conn.Close()
	}
	c.conn = conn
	c.conn.SetLogger(c.logger)

	c.logger.Printf("zkclient setup new connection to %s", c.zkAddr)

	go func() {
		for e := range events {
			log.Debugf("zookeeper event: %+v", e)
		}
	}()
	return nil
}

func (c *Client) Close() error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true

	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

func (c *Client) Do(fn func(conn *zk.Conn) error) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return errors.Trace(ErrClientClosed)
	}
	return c.shell(fn)
}

func (c *Client) shell(fn func(conn *zk.Conn) error) error {
	if err := fn(c.conn); err != nil {
		for _, e := range []error{zk.ErrNoNode, zk.ErrNodeExists, zk.ErrNotEmpty} {
			if errors.Equal(e, err) {
				return err
			}
		}
		if time.Since(c.dialAt) > time.Second {
			if err := c.reset(); err != nil {
				log.DebugErrorf(err, "zkclient reset connection failed")
			}
		}
		return err
	}
	return nil
}

func (c *Client) MkDir(path string) error {
	c.Lock()
	defer c.Unlock()

	if c.closed {
		return errors.Trace(ErrClientClosed)
	}

	log.Debugf("zkclient mkdir node %s", path)

	err := c.shell(func(conn *zk.Conn) error {
		return c.mkDir(conn, path)
	})

	if err != nil {
		log.Debugf("zkclient mkdir node %s failed: %s", path, err)
		return err
	}
	log.Debugf("zkclient mkdir OK")
	return nil
}

func (c *Client) mkDir(conn *zk.Conn, path string) error {
	if path == "" || path == "/" {
		return nil
	}
	if exists, _, err := conn.Exists(path); err != nil {
		return errors.Trace(err)
	} else if exists {
		return nil
	}
	if err := c.mkDir(conn, filepath.Dir(path)); err != nil {
		return err
	}
	_, err := conn.Create(path, []byte{}, 0, zk.WorldACL(zk.PermAll))
	if err != nil && errors.NotEqual(err, zk.ErrNodeExists) {
		return errors.Trace(err)
	}
	return nil
}

func (c *Client) Create(path string, data []byte) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return errors.Trace(ErrClientClosed)
	}
	log.Debugf("zkclient create node %s", path)
	err := c.shell(func(conn *zk.Conn) error {
		_, err := c.create(conn, path, data, 0)
		return err
	})
	if err != nil {
		log.Debugf("zkclient create node %s failed: %s", path, err)
		return err
	}
	log.Debugf("zkclient create OK")
	return nil
}

func (c *Client) create(conn *zk.Conn, path string, data []byte, flag int32) (string, error) {
	if err := c.mkDir(conn, filepath.Dir(path)); err != nil {
		return "", err
	}
	p, err := conn.Create(path, data, flag, zk.WorldACL(zk.PermAdmin|zk.PermRead|zk.PermWrite))
	if err != nil {
		return "", errors.Trace(err)
	}
	return p, nil
}

func (c *Client) watch(conn *zk.Conn, path string) (<-chan struct{}, error) {
	_, _, w, err := conn.GetW(path)
	if err != nil {
		return nil, errors.Trace(err)
	}
	signal := make(chan struct{})
	go func() {
		defer close(signal)
		<-w
		log.Debugf("zkclient watch node %s update", path)
	}()
	return signal, nil
}

func (c *Client) Update(path string, data []byte) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return errors.Trace(ErrClientClosed)
	}

	log.Debugf("zkclient update node %s", path)
	err := c.shell(func(conn *zk.Conn) error {
		return c.update(conn, path, data)
	})
	if err != nil {
		log.Debugf("zkclient update node %s failed: %s", path, err)
		return err
	}
	log.Debugf("zkclient update OK")
	return nil
}

func (c *Client) update(conn *zk.Conn, path string, data []byte) error {
	if exists, _, err := conn.Exists(path); err != nil {
		return errors.Trace(err)
	} else if !exists {
		_, err := c.create(conn, path, data, 0)
		if err != nil && errors.NotEqual(err, zk.ErrNodeExists) {
			return err
		}
	}
	_, err := conn.Set(path, data, -1)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (c *Client) Delete(path string) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return errors.Trace(ErrClientClosed)
	}

	log.Debugf("zkclient delete node %s", path)
	err := c.shell(func(conn *zk.Conn) error {
		err := conn.Delete(path, -1)
		if err != nil && errors.NotEqual(err, zk.ErrNoNode) {
			return errors.Trace(err)
		}
		return nil
	})
	if err != nil {
		log.Debugf("zkclient delete node %s failed: %s", path, err)
		return err
	}

	log.Debugf("zkclient delete OK")
	return nil
}

func (c *Client) Read(path string, must bool) ([]byte, error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil, errors.Trace(ErrClientClosed)
	}

	var data []byte
	err := c.shell(func(conn *zk.Conn) error {
		b, _, err := conn.Get(path)
		if err != nil {
			if errors.Equal(err, zk.ErrNoNode) && !must {
				return nil
			}
			return errors.Trace(err)
		}
		data = b
		return nil
	})

	if err != nil {
		log.Debugf("zkclient read node %s failed: %s", path, err)
		return nil, err
	}
	return data, nil
}


func (c *Client) List(path string, must bool) ([]string, error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil, errors.Trace(ErrClientClosed)
	}

	var paths []string
	err := c.shell(func(conn *zk.Conn) error {
		nodes, _, err := conn.Children(path)
		if err != nil {
			if errors.Equal(err, zk.ErrNoNode) && !must {
				return nil
			}
			return errors.Trace(err)
		}
		for _, node := range nodes {
			paths = append(paths, filepath.Join(path, node))
		}
		return nil
	})

	if err != nil {
		log.Debugf("zkclient list node %s failed: %s", path, err)
		return nil, err
	}
	return paths, nil
}
