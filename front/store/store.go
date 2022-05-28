package store

import (
	"github.com/ferealqq/wienerlist/front/actions"
	"github.com/ferealqq/wienerlist/front/dispatcher"
	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/ferealqq/wienerlist/front/store/services"
	"github.com/ferealqq/wienerlist/front/store/state"
	"github.com/ferealqq/wienerlist/front/store/storeutil"
)

var api = services.NewApi("http://localhost:4000/api/v1")

var (
	// Items represents all of the TODO items in the store.
	Items []*model.Item

	// Filter represents the active viewing filter for items.
	Filter = model.All

	// Listeners is the listeners that will be invoked when the store changes.
	Listeners = storeutil.NewListenerRegistry()

	SectionState = state.NewSectionState()
)

func init() {
	dispatcher.Register(onAction)
}

func FetchBoardSectionsIfNeeded(boardId int) error {
	if _, ok := SectionState.BoardSections[boardId]; !ok && !SectionState.IsFetching {
		dispatcher.Dispatch(&actions.FetchSectionsRequest{})
		var secs model.ListSections
		// TODO Create wrapper actions to get section data
		if err := api.Params("board_id", boardId).Get("/sections/").BindModel(&secs); err != nil {
			dispatcher.Dispatch(&actions.FetchSectionsResponseError{Error: err})
			return err
		}

		dispatcher.Dispatch(&actions.FetchSectionsResponse{
			Sections: secs.Sections,
			BoardId:  boardId,
		})
		return nil
	}

	return nil
}

// ActiveItemCount returns the current number of items that are not completed.
func ActiveItemCount() int {
	return count(false)
}

// CompletedItemCount returns the current number of items that are completed.
func CompletedItemCount() int {
	return count(true)
}

func count(completed bool) int {
	count := 0
	for _, item := range Items {
		if item.Completed == completed {
			count++
		}
	}
	return count
}

func onAction(action interface{}) {
	switch a := action.(type) {
	case *actions.ReplaceItems:
		Items = a.Items

	case *actions.AddItem:
		Items = append(Items, &model.Item{Title: a.Title, Completed: false})

	case *actions.DestroyItem:
		copy(Items[a.Index:], Items[a.Index+1:])
		Items = Items[:len(Items)-1]

	case *actions.SetTitle:
		Items[a.Index].Title = a.Title

	case *actions.SetCompleted:
		Items[a.Index].Completed = a.Completed

	case *actions.SetAllCompleted:
		for _, item := range Items {
			item.Completed = a.Completed
		}

	case *actions.ClearCompleted:
		var activeItems []*model.Item
		for _, item := range Items {
			if !item.Completed {
				activeItems = append(activeItems, item)
			}
		}
		Items = activeItems

	case *actions.SetFilter:
		Filter = a.Filter

	case *actions.FetchSectionsRequest:
		SectionState.IsFetching = true
		SectionState.LastActionFailed = false

	case *actions.FetchSectionsResponseError:
		SectionState.IsFetching = false
		SectionState.LastActionFailed = true
		SectionState.Error = a.Error

	case *actions.FetchSectionsResponse:
		l := len(a.Sections)
		// List of pointers
		secs := make([]*model.Section, 0, l)
		for i := 0; i != l; i++ {
			secs = append(secs, &a.Sections[i])
		}
		SectionState.BoardSections[a.BoardId] = secs
		SectionState.LastActionFailed = false
		SectionState.IsFetching = false

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}
