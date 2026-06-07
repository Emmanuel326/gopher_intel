# Gopher Intel

A Go bot that fetches activity from Linux kernel and systems software mailing lists,
then uses AI to produce structured intelligence briefs — so you can stay sharp on
what matters without drowning in list traffic.

## What it does

- Fetches recent messages from configured mailing lists (lore.kernel.org, SPDK, etc.)
- Batches all raw data into a single AI call (Gemini)
- Produces per-list intelligence briefs covering:
  - Situation overview
  - Key technical themes and threads
  - Notable patches and RFCs
  - Critical bugs and stability signals
  - People and review dynamics
  - Analyst take
  - Cross-list signals
- Falls back to a local digest (pure Go, no API) when Gemini quota is exhausted

## Sources

| List | URL |
|------|-----|
| Linux Kernel (LKML) | lore.kernel.org/lkml |
| Block Layer | lore.kernel.org/linux-block |
| NVMe | lore.kernel.org/linux-nvme |
| io_uring | lore.kernel.org/io-uring |
| QEMU / Virtio | lore.kernel.org/qemu-devel |
| BPF / eBPF | lore.kernel.org/bpf |

## Prerequisites

- Go 1.21+
- A Gemini API key (free tier) from https://aistudio.google.com/apikey

## Setup

```bash
git clone https://github.com/emmanuel326/gopher_intel
cd gopher_intel
cp .env.example .env
# add your key to .env
go mod download
```

## Configuration

Create a `.env` file in the project root:

```env
GEMINI_API_KEY=your_key_here
```

## Run

```bash
go run cmd/main.go
```

## Project structure
cat > .env.example << 'EOF'
GEMINI_API_KEY=your_key_here
