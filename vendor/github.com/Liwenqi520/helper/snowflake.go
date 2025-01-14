package helper

import "github.com/bwmarrin/snowflake"

//唯一自增主键
var snowflakeNode *snowflake.Node

type Snowflake struct{}

func (s Snowflake) Init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	snowflakeNode = node
}

// SetSnowflakeNode Set snowflake node
func (s Snowflake) SetSnowflakeNode(node, epoch int64) error {
	if epoch > 0 {
		snowflake.Epoch = epoch
	}
	n, err := snowflake.NewNode(node)
	if err != nil {
		return err
	}
	snowflakeNode = n
	return nil
}

// SnowflakeID Define alias
type SnowflakeID = snowflake.ID

// NewSnowflakeID Create snowflake id
func (s Snowflake) NewSnowflakeID() SnowflakeID {
	return snowflakeNode.Generate()
}

//NewID
func (s Snowflake) NewID() string {
	return s.NewSnowflakeID().String()
}
