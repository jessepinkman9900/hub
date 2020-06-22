package subscription

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hub "github.com/sentinel-official/hub/types"
	"github.com/sentinel-official/hub/x/dvpn/subscription/keeper"
	"github.com/sentinel-official/hub/x/dvpn/subscription/types"
)

func HandleAddPlan(ctx sdk.Context, k keeper.Keeper, msg types.MsgAddPlan) sdk.Result {
	_, found := k.GetProvider(ctx, msg.From)
	if !found {
		return types.ErrorNoProviderFound().Result()
	}

	count := k.GetPlansCountForProvider(ctx, msg.From)
	plan := types.Plan{
		ID:           count,
		Provider:     msg.From,
		Price:        msg.Price,
		Duration:     msg.Duration,
		MaxBandwidth: msg.MaxBandwidth,
		MaxDuration:  msg.MaxDuration,
		Status:       hub.StatusInactive,
		StatusAt:     ctx.BlockHeight(),
	}

	k.SetPlan(ctx, plan)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSetPlan,
		sdk.NewAttribute(types.AttributeKeyAddress, plan.Provider.String()),
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", plan.ID)),
	))

	k.SetPlansCountForProvider(ctx, plan.Provider, count+1)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSetPlansCount,
		sdk.NewAttribute(types.AttributeKeyAddress, plan.Provider.String()),
		sdk.NewAttribute(types.AttributeKeyCount, fmt.Sprintf("%d", count+1)),
	))

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleSetPlanStatus(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetPlanStatus) sdk.Result {
	plan, found := k.GetPlan(ctx, msg.From, msg.ID)
	if !found {
		return types.ErrorNoPlanFound().Result()
	}

	plan.Status = msg.Status
	plan.StatusAt = ctx.BlockHeight()

	k.SetPlan(ctx, plan)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSetPlanStatus,
		sdk.NewAttribute(types.AttributeKeyAddress, plan.Provider.String()),
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", plan.ID)),
		sdk.NewAttribute(types.AttributeKeyStatus, plan.Status.String()),
	))

	return sdk.Result{Events: ctx.EventManager().Events()}
}
