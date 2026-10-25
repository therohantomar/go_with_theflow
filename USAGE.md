# Usage

## 1. Run the CLI

```bash
cd go-with-the-flow
go mod tidy
go run ./cmd/flow new
```

## 2. (Optional) Install as a command

```bash
go build -o flow ./cmd/flow
sudo mv flow /usr/local/bin/
flow new
```

## 3. Run the generated site

```bash
cd your-project-name
npm install        # or: pnpm install / yarn / bun install
npm run dev        # http://localhost:4321
```

## 4. Build for production

```bash
npm run build      # output in dist/
npm run preview
```

## 5. Customize

Edit text: `src/content/<section>.json`

Add a photo:

```bash
cp ~/Pictures/factory.jpg src/assets/about/
```

then set it in `src/content/about.json`:

```json
"image": "factory.jpg"
```

Change the WhatsApp number (quote section) in `src/content/quote.json`:

```json
"whatsapp": "919999999999"
```

Change colors and fonts in `src/styles/global.css` (the `@theme` block).
