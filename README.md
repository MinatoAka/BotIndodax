<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.27.3%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.27.3+"></a>
  <a href="https://termux.dev/"><img src="https://img.shields.io/badge/Platform-Termux-2E8B57?style=for-the-badge&logo=android&logoColor=white" alt="Platform: Termux"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License: MIT"></a>
  <a href="https://saweria.co/AKaAturDompet"><img src="https://img.shields.io/badge/Support-Saweria-FF5722?style=for-the-badge&logo=buy-me-a-coffee&logoColor=white" alt="Support via Saweria"></a>
</p>

# BotIndodax

Automated Indodax trading bot written in Go for Termux. The bot monitors one IDR trading pair, evaluates a local buy/sell state, and submits market orders based on configured price triggers.

> [!CAUTION]
> This project can place real trades using your account balance. Trading carries financial risk and may result in losses. There is no guarantee of profit. Review the code and configuration before use. Use an API key with read and spot-trading permissions only; do not enable withdrawal access. Use this project at your own risk.

## Overview

BotIndodax is a lightweight command-line trading bot for Indodax designed to run on Android via Termux. It reads ticker data from the public market endpoint, loads account balances from the Indodax API, evaluates local trading state, and can submit BUY or SELL market orders depending on the configured trigger thresholds.

The bot is intentionally simple and focused on a single trading pair, with state persisted locally in `state.json`. It is intended for users who want a direct, script-driven trading workflow on mobile devices.

## Features

- Reads the latest ticker price for a configured pair such as `ethidr`.
- Retrieves free IDR and asset balances from the Indodax API.
- Uses a local `KOSONG` / `PEGANG` state to determine trade direction.
- Executes a BUY order when the price is at or below `buyTrigger` while the local state is `KOSONG`.
- Executes a SELL order when the price is at or above `sellTrigger` while the local state is `PEGANG`.
- Uses `modalPerTrade` to calculate the BUY quantity and checks `minQty` before sending an order.
- Applies `cooldownMs` to prevent repeated trades in quick succession.
- Reloads `config.json` on each iteration and persists local trading state to `state.json`.
- Stops execution when an order result is uncertain so users can reconcile with the Indodax account.

## Requirements

- Android device with Termux installed
- Internet connection
- Go 1.27.3 or newer
- Indodax account with TAPI v2 API access
- API key with read and spot trading permission only; withdrawal permission must remain disabled

## Installation

```sh
pkg update
pkg upgrade
pkg install git golang nano

git clone https://github.com/MinatoAka/BotIndodax.git
cd BotIndodax
```

## Configuration

Create or edit `config.json` in the repository root and populate your credentials and trading parameters.

```sh
nano config.json
```

Example configuration:

```json
{
  "apiKey": "ISI_API_KEY_LO",
  "secretKey": "ISI_SECRET_KEY_LO",
  "pair": "ethidr",
  "buyTrigger": 33450000,
  "sellTrigger": 34453500,
  "modalPerTrade": 15000,
  "minQty": 0.0001,
  "cooldownMs": 1800000
}
```

| Field | Description |
| --- | --- |
| `apiKey` | Indodax TAPI v2 API key. Required. |
| `secretKey` | Secret key paired with the API key. Required and must be kept private. |
| `pair` | Trading pair in lowercase, for example `ethidr`. |
| `buyTrigger` | Maximum price that triggers a BUY order. |
| `sellTrigger` | Minimum price that triggers a SELL order. Must be greater than `buyTrigger`. |
| `modalPerTrade` | IDR amount allocated for each BUY order. |
| `minQty` | Minimum asset quantity allowed for a BUY order. |
| `cooldownMs` | Minimum delay between trades, measured in milliseconds. |

The bot validates the configuration when it starts. If the file is invalid, it will log the error and continue retrying after one minute.

## Usage

Build the binary:

```sh
go build -o bot .
```

Run the bot:

```sh
./bot
```

Or run it directly without building:

```sh
go run .
```

For Android runtime stability, it is recommended to run the bot from the repository directory and, if needed, use:

```sh
termux-wake-lock
```

Stop the bot at any time with `Ctrl + C`.

## Safety Notes

- The bot stores local state in `state.json`.
- Default state is `KOSONG` if no file exists.
- If an order request times out or returns incomplete confirmation, do not immediately restart the bot. Reconcile the trade manually through the Indodax account first.
- Local state is not a substitute for checking balances and trade history in the actual exchange.
- Do not publish or commit API credentials in public repositories.

## Testing

The project includes lightweight unit tests for configuration validation and state handling.

```sh
go test ./...
```

Note: the current test suite references `MaxDailyTrades`, but that field is not present in the application configuration. This should be reviewed and aligned with the actual config model before the repository is treated as fully passing.

## Repository Structure

```text
BotIndodax/
├── LICENSE
├── README.md
├── config.go
├── config.json
├── go.mod
├── main.go
├── market.go
├── order.go
├── utils.go
├── utils_test.go
└── state.json   # generated at runtime
```

## Support

[♡ Support development via Saweria](https://saweria.co/AKaAturDompet)

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Disclaimer

This software is provided as-is, without warranty, and is not financial or investment advice. You are solely responsible for the security of your API keys, the correctness of your configuration, the monitoring of the bot, and the risks associated with automated trading.
