package core

import (
	"time"

	"bitbucket.org/ventureslash/go-ibft"
)

func (c *core) initTimeouts() {
	c.timeoutsMu.Lock()
	defer c.timeoutsMu.Unlock()
	c.logger.Info(c.address, ": Init timeouts for validators ")
	for _, val := range c.valSet.List() {
		if val.Address() != c.address {
			src := val.Address()
			c.logger.Info(c.address, ": Init timeout for ", src)
			c.timeouts[val] = time.AfterFunc(ibft.ValidatorTimeout, func() {
				c.logger.Info(c.address, ": Timeout:  deleting validator ", src)
				c.valSet.RemoveValidator(src)
			})
		}

	}
}

func (c *core) setValidatorTimeout(src ibft.Address) {
	c.timeoutsMu.Lock()
	defer c.timeoutsMu.Unlock()
	_, val := c.valSet.GetByAddress(src)
	if val != nil {
		c.logger.Info(c.address, ": Init timeout for ", src)
		if c.timeouts[val] != nil {
			c.timeouts[val].Stop()
		}
		src := val.Address()

		c.timeouts[val] = time.AfterFunc(ibft.ValidatorTimeout, func() {

			c.logger.Info(c.address, ": Timeout: deleting validator ", src)
			c.valSet.RemoveValidator(src)
		})
	}
}
