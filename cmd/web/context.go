package main

type contextKey string

// Used to obtain context var from user request
const isAuthenticatedContextKey = contextKey("isAuthenticated")
