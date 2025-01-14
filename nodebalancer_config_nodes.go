package linodego

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// NodeBalancerNode objects represent a backend that can accept traffic for a NodeBalancer Config
type NodeBalancerNode struct {
	ID             int      `json:"id"`
	Address        string   `json:"address"`
	Label          string   `json:"label"`
	Status         string   `json:"status"`
	Weight         int      `json:"weight"`
	Mode           NodeMode `json:"mode"`
	ConfigID       int      `json:"config_id"`
	NodeBalancerID int      `json:"nodebalancer_id"`
}

// NodeMode is the mode a NodeBalancer should use when sending traffic to a NodeBalancer Node
type NodeMode string

var (
	// ModeAccept is the NodeMode indicating a NodeBalancer Node is accepting traffic
	ModeAccept NodeMode = "accept"

	// ModeReject is the NodeMode indicating a NodeBalancer Node is not receiving traffic
	ModeReject NodeMode = "reject"

	// ModeDrain is the NodeMode indicating a NodeBalancer Node is not receiving new traffic, but may continue receiving traffic from pinned connections
	ModeDrain NodeMode = "drain"

	// ModeBackup is the NodeMode indicating a NodeBalancer Node will only receive traffic if all "accept" Nodes are down
	ModeBackup NodeMode = "backup"
)

// NodeBalancerNodeCreateOptions fields are those accepted by CreateNodeBalancerNode
type NodeBalancerNodeCreateOptions struct {
	Address string   `json:"address"`
	Label   string   `json:"label"`
	Weight  int      `json:"weight,omitempty"`
	Mode    NodeMode `json:"mode,omitempty"`
}

// NodeBalancerNodeUpdateOptions fields are those accepted by UpdateNodeBalancerNode
type NodeBalancerNodeUpdateOptions struct {
	Address string   `json:"address,omitempty"`
	Label   string   `json:"label,omitempty"`
	Weight  int      `json:"weight,omitempty"`
	Mode    NodeMode `json:"mode,omitempty"`
}

// GetCreateOptions converts a NodeBalancerNode to NodeBalancerNodeCreateOptions for use in CreateNodeBalancerNode
func (i NodeBalancerNode) GetCreateOptions() NodeBalancerNodeCreateOptions {
	return NodeBalancerNodeCreateOptions{
		Address: i.Address,
		Label:   i.Label,
		Weight:  i.Weight,
		Mode:    i.Mode,
	}
}

// GetUpdateOptions converts a NodeBalancerNode to NodeBalancerNodeUpdateOptions for use in UpdateNodeBalancerNode
func (i NodeBalancerNode) GetUpdateOptions() NodeBalancerNodeUpdateOptions {
	return NodeBalancerNodeUpdateOptions{
		Address: i.Address,
		Label:   i.Label,
		Weight:  i.Weight,
		Mode:    i.Mode,
	}
}

// NodeBalancerNodesPagedResponse represents a paginated NodeBalancerNode API response
type NodeBalancerNodesPagedResponse struct {
	*PageOptions
	Data []NodeBalancerNode `json:"data"`
}

// endpoint gets the endpoint URL for NodeBalancerNode
func (NodeBalancerNodesPagedResponse) endpoint(c *Client, ids ...any) string {
	nodebalancerID := ids[0].(int)
	configID := ids[1].(int)
	endpoint, err := c.NodeBalancerNodes.endpointWithParams(nodebalancerID, configID)
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *NodeBalancerNodesPagedResponse) castResult(r *resty.Request, e string) (int, int, error) {
	res, err := coupleAPIErrors(r.SetResult(NodeBalancerNodesPagedResponse{}).Get(e))
	if err != nil {
		return 0, 0, err
	}
	castedRes := res.Result().(*NodeBalancerNodesPagedResponse)
	resp.Data = append(resp.Data, castedRes.Data...)
	return castedRes.Pages, castedRes.Results, nil
}

// ListNodeBalancerNodes lists NodeBalancerNodes
func (c *Client) ListNodeBalancerNodes(ctx context.Context, nodebalancerID int, configID int, opts *ListOptions) ([]NodeBalancerNode, error) {
	response := NodeBalancerNodesPagedResponse{}
	err := c.listHelper(ctx, &response, opts, nodebalancerID, configID)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetNodeBalancerNode gets the template with the provided ID
func (c *Client) GetNodeBalancerNode(ctx context.Context, nodebalancerID int, configID int, nodeID int) (*NodeBalancerNode, error) {
	e, err := c.NodeBalancerNodes.endpointWithParams(nodebalancerID, configID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, nodeID)
	r, err := coupleAPIErrors(c.R(ctx).SetResult(&NodeBalancerNode{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancerNode), nil
}

// CreateNodeBalancerNode creates a NodeBalancerNode
func (c *Client) CreateNodeBalancerNode(ctx context.Context, nodebalancerID int, configID int, createOpts NodeBalancerNodeCreateOptions) (*NodeBalancerNode, error) {
	var body string
	e, err := c.NodeBalancerNodes.endpointWithParams(nodebalancerID, configID)
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&NodeBalancerNode{})

	if bodyData, err := json.Marshal(createOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Post(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancerNode), nil
}

// UpdateNodeBalancerNode updates the NodeBalancerNode with the specified id
func (c *Client) UpdateNodeBalancerNode(ctx context.Context, nodebalancerID int, configID int, nodeID int, updateOpts NodeBalancerNodeUpdateOptions) (*NodeBalancerNode, error) {
	var body string
	e, err := c.NodeBalancerNodes.endpointWithParams(nodebalancerID, configID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, nodeID)

	req := c.R(ctx).SetResult(&NodeBalancerNode{})

	if bodyData, err := json.Marshal(updateOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Put(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancerNode), nil
}

// DeleteNodeBalancerNode deletes the NodeBalancerNode with the specified id
func (c *Client) DeleteNodeBalancerNode(ctx context.Context, nodebalancerID int, configID int, nodeID int) error {
	e, err := c.NodeBalancerNodes.endpointWithParams(nodebalancerID, configID)
	if err != nil {
		return err
	}
	e = fmt.Sprintf("%s/%d", e, nodeID)

	_, err = coupleAPIErrors(c.R(ctx).Delete(e))
	return err
}
