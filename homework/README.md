### How to Use the Project

#### For Windows Users
If you are using a Windows operating system, please follow these steps:

1. Open your command line tool, such as CMD or PowerShell.
2. Navigate to the `homework` directory path.
3. Build the project by executing the following command:
   ```
   go build
   ```
4. Run the generated executable file:
   ```
   .\homework.exe
   ```
5. After execution, you will be able to see the current block and all transactions within the block, including the hash, sender (from), recipient (to), and transaction value (value).

### Adding Subscribers

To add subscribers, you can manually call the following methods within the `main` function:

- `EthereumParser.Subscribe`: Subscribe to a specific Ethereum address.
- `EthereumParser.GetTransactions`: Retrieve all transactions for the subscribed address.

Using these methods, you can dynamically obtain and track all transaction activities for the subscribed address.