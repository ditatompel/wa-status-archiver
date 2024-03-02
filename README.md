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

And many indirect dependencies can be found in `go.mod`, `go.sum`, `package.json` and `package-lock.json`.

## Development

```shell
npm install
make copyhtmx
npm run dev
```

## Build binary file

```shell
npm ci
make build
```

