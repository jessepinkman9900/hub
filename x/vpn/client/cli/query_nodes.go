package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	csdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdkTypes "github.com/ironman0x7b2/sentinel-sdk/types"
	"github.com/ironman0x7b2/sentinel-sdk/x/vpn"
	"github.com/ironman0x7b2/sentinel-sdk/x/vpn/client/common"
)

func QueryNodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Get node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			res, err := common.QueryNode(cliCtx, cdc, sdkTypes.NewID(args[0]))
			if err != nil {
				return err
			}

			nodeData, err := cdc.MarshalJSONIndent(res, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(nodeData))

			return nil
		},
	}

	return cmd
}

func QueryNodesCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Get nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			owner := viper.GetString(flagOwnerAddress)

			var nodes []vpn.Node
			if len(owner) == 0 {
				res, err := cliCtx.QuerySubspace(vpn.NodeKeyPrefix, vpn.StoreKeyNode)
				if err != nil {
					return err
				}
				if len(res) == 0 {
					return fmt.Errorf("no nodes found")
				}

				for _, kv := range res {
					var node vpn.Node
					if err := cdc.UnmarshalBinaryLengthPrefixed(kv.Value, &node); err != nil {
						return err
					}

					nodes = append(nodes, node)
				}
			} else {
				owner, err := csdkTypes.AccAddressFromBech32(owner)
				if err != nil {
					return err
				}

				res, err := common.QueryNodesOfOwner(cliCtx, cdc, owner)
				if err != nil {
					return err
				}
				if string(res) == "null" {
					return fmt.Errorf("no nodes found")
				}

				if err := cdc.UnmarshalJSON(res, &nodes); err != nil {
					return err
				}
			}

			nodesData, err := cdc.MarshalJSONIndent(nodes, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(nodesData))

			return nil
		},
	}

	cmd.Flags().String(flagOwnerAddress, "", "VPN node owner address")

	return cmd
}
