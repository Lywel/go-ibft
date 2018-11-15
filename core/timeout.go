package core

import (
	"time"

	"bitbucket.org/ventureslash/go-ibft"
)

func (c *core) setValidatorTimeout(src ibft.Address) {
	c.timeoutsMu.Lock()
	defer c.timeoutsMu.Unlock()
	_, val := c.valSet.GetByAddress(src)
	if val != nil {
		c.timeouts[val] = time.AfterFunc(ibft.ValidatorTimeout, func() {
			c.valSet.RemoveValidator(val.Address())
		})
	}
}
