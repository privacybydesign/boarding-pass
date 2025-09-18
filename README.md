# Boarding Pass

This is a demonstration of one of the usecases of using Yivi passport
credentials. By obtaining the passport credential and only disclosing certain
attributes you can check-in to your flight and receive a boarding pass.

To run everything locally with Docker, place your jwt keys in `local-secrets/`,
then run `docker compose build` and `docker compose up`. The compose stack

for running locally without docker, you need to first build the frontend with
`npm run build`. The go server will host the frontend. Then do `go build .` in
the backend directory and the start the server with
`go run . --config config.json` from the `backend/` directory and launch the
frontend with `npm install && npm run dev` inside `frontend/`; point
`frontend/.env` at the backend URL you just started.
