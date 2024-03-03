# WA Status Archiver

> This is **experimental**, and **DO NOT** expose the WEB UI to the public!

## Attribution

Attribution comes first, this project is build on top of these following awesome open source projects:

-   [go.mau.fi/whatsmeow](https://github.com/tulir/whatsmeow/)
-   [HTMX](https://htmx.org/)
-   [Tailwind CSS](https://tailwindcss.com/)
-   [gofiber/fiber](https://github.com/gofiber/fiber)
-   [mdp/qrterminal](https://github.com/mdp/qrterminal)
-   [spf13/cobra](https://github.com/spf13/cobra)
-   [joho/godotenv](https://github.com/joho/godotenv)
-   [alexedwards/argon2id](https://github.com/alexedwards/argon2id)
-   [golang.org/x/term](https://pkg.go.dev/golang.org/x/term)
-   [gosimple/slug](https://github.com/gosimple/slug)
-   [Flowbite](https://flowbite.com) (Tailwind design components)

And many indirect dependencies can be found in `go.mod`, `go.sum`, `package.json` and `package-lock.json`.

## Trying this app

### Requirements (tested on)

-   Node.js >=20.x
-   Go 1.22.x
-   PostgreSQL >=15.x

### Prepare the assets

> Only run these steps once.

1. Clone this repository
2. Copy `.env.example` to `.env` and modify as needed (especially `SECRET_KEY` and **DB** config)
3. run `npm ci`
4. run `make static`

### Running the "bot"

1. Run `go run . run`
2. On initial setup, it will display QR code. Scan it by linking device with your phone.

### Running the web UI

To access Web UI, you need to create an admin account. This can be done by running `go run . admin create` and fill your username and password.

After that, you can run `go run . serve` and access the UI from the browser. (default: http://127.0.0.1:18088)

## Build the binary file

```shell
make build
```

The binary file will be placed in the `bin` directory.

## Development

If you want to develop or modify UI:

```shell
npm install
make static
npm run dev
go run . serve
```

Thanks to [cosmtrek/air](https://github.com/cosmtrek/air), you can run `air serve` to live reload the HTTP server (Do not use `air` when running `run` command).

## FAQ

### Is this project stable?

Definitely **no**.

### What database is supported?

Although `whatsmeow` support SQLite and PostgreSQL, I only create this project on top of PostgreSQL. Feel free to adapt the database driver to fit with your needs.

### Is this project support for multiple account?

No, this project only support one account.

