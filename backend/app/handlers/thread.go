package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"tech-db-forum/app/models"
	"tech-db-forum/pkg/database"
	"tech-db-forum/pkg/network"

	"github.com/gorilla/mux"
)

func ThreadCreate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	thread := models.Thread{}
	_ = thread.UnmarshalJSON(body)

	user, err := database.GetUserByNickname(thread.Author)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	thread.Author = user.Nickname

	forumSlug := mux.Vars(r)["slug"]
	forum, err := database.GetForumBySlug(forumSlug)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	thread.Forum = forum.Slug

	if thread.Slug != "" {
		existingThread, err := database.GetThreadBySlug(thread.Slug)
		if err == nil {
			network.WriteResponse(w, http.StatusConflict, existingThread)
			return
		}
	}

	err = database.CreateThread(&thread)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusCreated, thread)
}

func ThreadGetOne(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]
	thread, err := database.GetThread(slugOrId)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	network.WriteResponse(w, http.StatusOK, thread)
}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]
	thread, err := database.GetThread(slugOrId)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	args := r.URL.Query()
	limit := args.Get("limit")
	if limit == "" {
		limit = "1"
	}
	since := args.Get("since")
	sort := args.Get("sort")
	if sort == "" {
		sort = "flat"
	}
	desc, _ := strconv.ParseBool(args.Get("desc"))

	posts, err := database.GetThreadPosts(thread.Id, limit, since, sort, desc)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusOK, posts)
}

func ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]
	thread, err := database.GetThread(slugOrId)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	threadToUpdate := models.Thread{}
	_ = threadToUpdate.UnmarshalJSON(body)
	threadToUpdate.Slug = thread.Slug
	log.Print(threadToUpdate)

	err = database.UpdateThread(&threadToUpdate)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusOK, threadToUpdate)
}

func ThreadVote(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]
	thread, err := database.GetThread(slugOrId)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	vote := models.Vote{}
	_ = vote.UnmarshalJSON(body)

	_, err = database.GetUserByNickname(vote.Nickname)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	newVotes, err := database.Vote(vote, thread.Id)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	thread.Votes = newVotes

	network.WriteResponse(w, http.StatusOK, thread)
}
