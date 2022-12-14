package handlers

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"tech-db-forum/app/models"
	"tech-db-forum/pkg/database"
	"tech-db-forum/pkg/network"

	"github.com/gorilla/mux"
)

func ForumCreate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	forum := models.Forum{}
	_ = forum.UnmarshalJSON(body)

	user, err := database.GetUserByNickname(forum.User)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	forum.User = user.Nickname

	existingForum, err := database.GetForumBySlug(forum.Slug)
	if err == nil {
		network.WriteResponse(w, http.StatusConflict, existingForum)
		return
	}

	err = database.CreateForum(forum)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusCreated, forum)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	forum, err := database.GetForumBySlug(slug)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusOK, forum)
}

func ForumGetThreads(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	forum, err := database.GetForumBySlug(slug)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	slug = forum.Slug

	args := r.URL.Query()
	limit := args.Get("limit")
	if limit == "" {
		limit = "1"
	}
	since := args.Get("since")
	desc, _ := strconv.ParseBool(args.Get("desc"))

	threads, err := database.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusOK, threads)
}

func ForumGetUsers(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	forum, err := database.GetForumBySlug(slug)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	slug = forum.Slug

	args := r.URL.Query()
	limit := args.Get("limit")
	if limit == "" {
		limit = "1"
	}
	since := args.Get("since")
	desc, _ := strconv.ParseBool(args.Get("desc"))

	users, err := database.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusOK, users)
}
