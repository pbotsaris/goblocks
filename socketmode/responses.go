package socketmode

import (
	"github.com/pbotsaris/goblocks/blocks"
)

// Response is the interface for all handler responses.
// The unexported method seals the interface to this package.
type Response interface {
	response()
	// toPayload converts the response to the ack payload format.
	toPayload() any
}

// ModalAction specifies what to do with a modal response.
type ModalAction int

const (
	// ModalActionUpdate replaces the current modal view.
	ModalActionUpdate ModalAction = iota
	// ModalActionPush pushes a new modal onto the stack.
	ModalActionPush
	// ModalActionClear closes all modal views.
	ModalActionClear
	// ModalActionErrors returns validation errors.
	ModalActionErrors
)

// EmptyResponse represents an acknowledgment with no payload.
type EmptyResponse struct{}

func (EmptyResponse) response()      {}
func (EmptyResponse) toPayload() any { return nil }

// MessageResponse responds with a message (for slash commands).
type MessageResponse struct {
	message blocks.Message
}

func (MessageResponse) response() {}
func (r MessageResponse) toPayload() any {
	return r.message
}

// ModalResponse responds with a modal action (for view submissions).
type ModalResponse struct {
	action ModalAction
	view   *blocks.Modal
	errors map[string]string
}

func (ModalResponse) response() {}
func (r ModalResponse) toPayload() any {
	switch r.action {
	case ModalActionUpdate:
		if r.view == nil {
			return nil
		}
		return map[string]any{
			"response_action": "update",
			"view":            r.view,
		}
	case ModalActionPush:
		if r.view == nil {
			return nil
		}
		return map[string]any{
			"response_action": "push",
			"view":            r.view,
		}
	case ModalActionClear:
		return map[string]any{
			"response_action": "clear",
		}
	case ModalActionErrors:
		if len(r.errors) == 0 {
			return nil
		}
		return map[string]any{
			"response_action": "errors",
			"errors":          r.errors,
		}
	default:
		return nil
	}
}

// OptionsResponse responds with options (for block_suggestion/dynamic menus).
type OptionsResponse struct {
	options      []blocks.Option
	optionGroups []blocks.OptionGroup
}

func (OptionsResponse) response() {}
func (r OptionsResponse) toPayload() any {
	if len(r.optionGroups) > 0 {
		return map[string]any{
			"option_groups": r.optionGroups,
		}
	}
	return map[string]any{
		"options": r.options,
	}
}

// Response builders

// NoResponse returns an empty response (ack only, no payload).
func NoResponse() Response {
	return EmptyResponse{}
}

// RespondWithMessage creates a response with a message.
// Use for slash commands that accept response payloads.
func RespondWithMessage(msg blocks.Message) Response {
	return MessageResponse{message: msg}
}

// RespondWithBlocks creates a message response from blocks directly.
// Convenience wrapper around RespondWithMessage.
func RespondWithBlocks(blks []blocks.Block) Response {
	msg, err := blocks.NewMessageWithBlocks(blks)
	if err != nil {
		// If blocks are invalid, return empty response
		return EmptyResponse{}
	}
	return MessageResponse{message: msg}
}

// RespondWithModalUpdate creates a response that updates the current modal.
// Use for view_submission when you want to replace the modal content.
func RespondWithModalUpdate(modal blocks.Modal) Response {
	return ModalResponse{
		action: ModalActionUpdate,
		view:   &modal,
	}
}

// RespondWithModalPush creates a response that pushes a new modal.
// Use for view_submission when you want to add a modal to the stack.
func RespondWithModalPush(modal blocks.Modal) Response {
	return ModalResponse{
		action: ModalActionPush,
		view:   &modal,
	}
}

// RespondWithModalClear creates a response that closes all modals.
// Use for view_submission when you want to dismiss the modal stack.
func RespondWithModalClear() Response {
	return ModalResponse{
		action: ModalActionClear,
	}
}

// RespondWithErrors creates a response with validation errors.
// Use for view_submission when form validation fails.
// The errors map keys are block_ids, values are error messages.
func RespondWithErrors(errors map[string]string) Response {
	return ModalResponse{
		action: ModalActionErrors,
		errors: errors,
	}
}

// RespondWithOptions creates a response with options for dynamic menus.
// Use for block_suggestion when populating an external select.
func RespondWithOptions(opts []blocks.Option) Response {
	return OptionsResponse{options: opts}
}

// RespondWithOptionGroups creates a response with grouped options.
// Use for block_suggestion when populating an external select with groups.
func RespondWithOptionGroups(groups []blocks.OptionGroup) Response {
	return OptionsResponse{optionGroups: groups}
}
