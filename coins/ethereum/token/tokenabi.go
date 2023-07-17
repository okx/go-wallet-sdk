package token

const ERC20ABI = `[
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_spender",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "approve",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "totalSupply",
		"outputs": [
			{
				"name": "supply",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_from",
				"type": "address"
			},
			{
				"name": "_to",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "transferFrom",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "decimals",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "_owner",
				"type": "address"
			}
		],
		"name": "balanceOf",
		"outputs": [
			{
				"name": "balance",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "symbol",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_to",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "transfer",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "_owner",
				"type": "address"
			},
			{
				"name": "_spender",
				"type": "address"
			}
		],
		"name": "allowance",
		"outputs": [
			{
				"name": "remaining",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "_from",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "_to",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "Transfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "_owner",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "_spender",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "Approval",
		"type": "event"
	}
]`
const ERC721ABI = `[
    {
      "constant": true,
      "inputs": [
        {
          "name": "interfaceId",
          "type": "bytes4"
        }
      ],
      "name": "supportsInterface",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "getApproved",
      "outputs": [
        {
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "to",
          "type": "address"
        },
        {
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "from",
          "type": "address"
        },
        {
          "name": "to",
          "type": "address"
        },
        {
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "transferFrom",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "from",
          "type": "address"
        },
        {
          "name": "to",
          "type": "address"
        },
        {
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "safeTransferFrom",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "ownerOf",
      "outputs": [
        {
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "to",
          "type": "address"
        },
        {
          "name": "approved",
          "type": "bool"
        }
      ],
      "name": "setApprovalForAll",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "owner",
          "type": "address"
        },
        {
          "name": "operator",
          "type": "address"
        }
      ],
      "name": "isApprovedForAll",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "name": "name",
          "type": "string"
        },
        {
          "name": "symbol",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "name": "to",
          "type": "address"
        },
        {
          "indexed": true,
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "name": "approved",
          "type": "address"
        },
        {
          "indexed": true,
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "Approval",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "name": "operator",
          "type": "address"
        },
        {
          "indexed": false,
          "name": "approved",
          "type": "bool"
        }
      ],
      "name": "ApprovalForAll",
      "type": "event"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "name",
      "outputs": [
        {
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "symbol",
      "outputs": [
        {
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "tokenURI",
      "outputs": [
        {
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    }
  ]`

const ERC1155ABI = `[
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "orderId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "operator",
        "type": "address"
      }
    ],
    "name": "CancelOrder",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "orderId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "maker",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "enum FixedPriceTrade1155.OrderTypeV2",
        "name": "orderType",
        "type": "uint8"
      }
    ],
    "name": "CancelOrderV2",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "address",
        "name": "previousAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "currentAddress",
        "type": "address"
      }
    ],
    "name": "ERC1155AddressWithCopyrightChanged",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "orderId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "maker",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "tokenAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "id",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "remainingAmount",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "payTokenAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "fixedPrice",
        "type": "uint256"
      }
    ],
    "name": "MakeOrder",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "orderId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "maker",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "enum FixedPriceTrade1155.OrderTypeV2",
        "name": "orderType",
        "type": "uint8"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "payTokenAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "price",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "tokenAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "id",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "expiredTime",
        "type": "uint256"
      }
    ],
    "name": "MakeOrderV2",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "previousOwner",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "newOwner",
        "type": "address"
      }
    ],
    "name": "OwnershipTransferred",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "orderId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "taker",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "author",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "copyrightFee",
        "type": "uint256"
      }
    ],
    "name": "PayCopyrightFee",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "address",
        "name": "erc20Addr",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "bool",
        "name": "jurisdiction",
        "type": "bool"
      }
    ],
    "name": "PaymentWhitelistChange",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "orderId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "taker",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "maker",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "tokenAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "id",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "payTokenAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "fixedPrice",
        "type": "uint256"
      }
    ],
    "name": "TakeOrder",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "orderId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "receiverERC1155",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "receiverERC20",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "enum FixedPriceTrade1155.OrderTypeV2",
        "name": "orderType",
        "type": "uint8"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "tokenAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "id",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "payTokenAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "price",
        "type": "uint256"
      }
    ],
    "name": "TakeOrderV2",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "orderId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "operator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "newRemainingAmount",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "newPayTokenAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "newFixedPrice",
        "type": "uint256"
      }
    ],
    "name": "UpdateOrder",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "orderId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "maker",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "enum FixedPriceTrade1155.OrderTypeV2",
        "name": "orderType",
        "type": "uint8"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "newAmount",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "newPayTokenAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "newPrice",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "newExpiredTime",
        "type": "uint256"
      }
    ],
    "name": "UpdateOrderV2",
    "type": "event"
  },
  {
    "inputs": [
      { "internalType": "address", "name": "tokenAddress", "type": "address" },
      { "internalType": "uint256", "name": "id", "type": "uint256" },
      { "internalType": "uint256", "name": "amount", "type": "uint256" },
      {
        "internalType": "address",
        "name": "payTokenAddress",
        "type": "address"
      },
      { "internalType": "uint256", "name": "fixedPrice", "type": "uint256" }
    ],
    "name": "ask",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "uint256", "name": "orderId", "type": "uint256" },
      { "internalType": "uint256", "name": "amount", "type": "uint256" }
    ],
    "name": "bid",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "uint256", "name": "orderId", "type": "uint256" }
    ],
    "name": "cancelOrder",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "uint256", "name": "orderID", "type": "uint256" }
    ],
    "name": "cancelOrderV2",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "getERC1155AddressWithCopyright",
    "outputs": [{ "internalType": "address", "name": "", "type": "address" }],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "uint256", "name": "orderId", "type": "uint256" }
    ],
    "name": "getOrder",
    "outputs": [
      {
        "components": [
          { "internalType": "address", "name": "maker", "type": "address" },
          {
            "internalType": "address",
            "name": "tokenAddress",
            "type": "address"
          },
          { "internalType": "uint256", "name": "id", "type": "uint256" },
          {
            "internalType": "uint256",
            "name": "remainingAmount",
            "type": "uint256"
          },
          {
            "internalType": "address",
            "name": "payTokenAddress",
            "type": "address"
          },
          {
            "internalType": "uint256",
            "name": "fixedPrice",
            "type": "uint256"
          },
          { "internalType": "bool", "name": "isAvailable", "type": "bool" }
        ],
        "internalType": "struct FixedPriceTrade1155.Order",
        "name": "order",
        "type": "tuple"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "uint256", "name": "orderId", "type": "uint256" }
    ],
    "name": "getOrderV2",
    "outputs": [
      {
        "components": [
          {
            "internalType": "enum FixedPriceTrade1155.OrderTypeV2",
            "name": "orderType",
            "type": "uint8"
          },
          {
            "components": [
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "tokenAddress",
                "type": "address"
              },
              { "internalType": "uint256", "name": "id", "type": "uint256" },
              {
                "internalType": "uint256",
                "name": "remainingAmount",
                "type": "uint256"
              },
              {
                "internalType": "address",
                "name": "payTokenAddress",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "fixedPrice",
                "type": "uint256"
              },
              { "internalType": "bool", "name": "isAvailable", "type": "bool" }
            ],
            "internalType": "struct FixedPriceTrade1155.Order",
            "name": "order",
            "type": "tuple"
          },
          {
            "internalType": "uint256",
            "name": "expiredTime",
            "type": "uint256"
          }
        ],
        "internalType": "struct FixedPriceTrade1155.OrderV2",
        "name": "orderV2",
        "type": "tuple"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "address", "name": "erc20Addr", "type": "address" }
    ],
    "name": "getPaymentWhitelist",
    "outputs": [{ "internalType": "bool", "name": "", "type": "bool" }],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "ERC1155AddressWithCopyright",
        "type": "address"
      },
      { "internalType": "address", "name": "newOwner", "type": "address" }
    ],
    "name": "init",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "payTokenAddress",
        "type": "address"
      },
      { "internalType": "uint256", "name": "price", "type": "uint256" },
      { "internalType": "address", "name": "tokenAddress", "type": "address" },
      { "internalType": "uint256", "name": "id", "type": "uint256" },
      { "internalType": "uint256", "name": "amount", "type": "uint256" },
      { "internalType": "uint256", "name": "expireTime", "type": "uint256" }
    ],
    "name": "makeOfferV2",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "address", "name": "tokenAddress", "type": "address" },
      { "internalType": "uint256", "name": "id", "type": "uint256" },
      { "internalType": "uint256", "name": "amount", "type": "uint256" },
      {
        "internalType": "address",
        "name": "payTokenAddress",
        "type": "address"
      },
      { "internalType": "uint256", "name": "price", "type": "uint256" },
      { "internalType": "uint256", "name": "expireTime", "type": "uint256" }
    ],
    "name": "makeSaleV2",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "owner",
    "outputs": [{ "internalType": "address", "name": "", "type": "address" }],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "renounceOwnership",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "newERC1155AddressWithCopyright",
        "type": "address"
      }
    ],
    "name": "setERC1155AddressWithCopyright",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "address", "name": "erc20Addr", "type": "address" },
      { "internalType": "bool", "name": "jurisdiction", "type": "bool" }
    ],
    "name": "setPaymentWhitelist",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "uint256", "name": "orderId", "type": "uint256" },
      { "internalType": "uint256", "name": "amount", "type": "uint256" }
    ],
    "name": "takeOrderV2",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "address", "name": "newOwner", "type": "address" }
    ],
    "name": "transferOwnership",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "uint256", "name": "orderId", "type": "uint256" },
      { "internalType": "uint256", "name": "newAmount", "type": "uint256" },
      {
        "internalType": "address",
        "name": "newPayTokenAddress",
        "type": "address"
      },
      { "internalType": "uint256", "name": "newFixedPrice", "type": "uint256" }
    ],
    "name": "updateOrder",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "internalType": "uint256", "name": "orderId", "type": "uint256" },
      { "internalType": "uint256", "name": "newAmount", "type": "uint256" },
      {
        "internalType": "address",
        "name": "newPayTokenAddress",
        "type": "address"
      },
      { "internalType": "uint256", "name": "newPrice", "type": "uint256" },
      { "internalType": "uint256", "name": "newExpiredTime", "type": "uint256" }
    ],
    "name": "updateOrderV2",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  }
]
`

const ERCBytes1155ABI = `
[
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "orderId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "operator",
        "type": "address"
      }
    ],
    "name": "CancelOrder",
    "type": "event"
  },
{
      "inputs": [
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "id",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "data",
          "type": "bytes"
        }
      ],
      "name": "safeTransferFrom",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
]
`

const AddressArrayABI = `[
  {
    "constant": true,
    "inputs": [],
    "name": "name",
    "outputs": [{ "name": "", "type": "string" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "tokenTransferProxy",
    "outputs": [{ "name": "", "type": "address" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "target", "type": "address" },
      { "name": "calldata", "type": "bytes" },
      { "name": "extradata", "type": "bytes" }
    ],
    "name": "staticCall",
    "outputs": [{ "name": "result", "type": "bool" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [{ "name": "newMinimumMakerProtocolFee", "type": "uint256" }],
    "name": "changeMinimumMakerProtocolFee",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [{ "name": "newMinimumTakerProtocolFee", "type": "uint256" }],
    "name": "changeMinimumTakerProtocolFee",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "array", "type": "bytes" },
      { "name": "desired", "type": "bytes" },
      { "name": "mask", "type": "bytes" }
    ],
    "name": "guardedArrayReplace",
    "outputs": [{ "name": "", "type": "bytes" }],
    "payable": false,
    "stateMutability": "pure",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "minimumTakerProtocolFee",
    "outputs": [{ "name": "", "type": "uint256" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "codename",
    "outputs": [{ "name": "", "type": "string" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [{ "name": "addr", "type": "address" }],
    "name": "testCopyAddress",
    "outputs": [{ "name": "", "type": "bytes" }],
    "payable": false,
    "stateMutability": "pure",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [{ "name": "arrToCopy", "type": "bytes" }],
    "name": "testCopy",
    "outputs": [{ "name": "", "type": "bytes" }],
    "payable": false,
    "stateMutability": "pure",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "addrs", "type": "address[7]" },
      { "name": "uints", "type": "uint256[9]" },
      { "name": "feeMethod", "type": "uint8" },
      { "name": "side", "type": "uint8" },
      { "name": "saleKind", "type": "uint8" },
      { "name": "howToCall", "type": "uint8" },
      { "name": "calldata", "type": "bytes" },
      { "name": "replacementPattern", "type": "bytes" },
      { "name": "staticExtradata", "type": "bytes" }
    ],
    "name": "calculateCurrentPrice_",
    "outputs": [{ "name": "", "type": "uint256" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [{ "name": "newProtocolFeeRecipient", "type": "address" }],
    "name": "changeProtocolFeeRecipient",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "version",
    "outputs": [{ "name": "", "type": "string" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "buyCalldata", "type": "bytes" },
      { "name": "buyReplacementPattern", "type": "bytes" },
      { "name": "sellCalldata", "type": "bytes" },
      { "name": "sellReplacementPattern", "type": "bytes" }
    ],
    "name": "orderCalldataCanMatch",
    "outputs": [{ "name": "", "type": "bool" }],
    "payable": false,
    "stateMutability": "pure",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "addrs", "type": "address[7]" },
      { "name": "uints", "type": "uint256[9]" },
      { "name": "feeMethod", "type": "uint8" },
      { "name": "side", "type": "uint8" },
      { "name": "saleKind", "type": "uint8" },
      { "name": "howToCall", "type": "uint8" },
      { "name": "calldata", "type": "bytes" },
      { "name": "replacementPattern", "type": "bytes" },
      { "name": "staticExtradata", "type": "bytes" },
      { "name": "v", "type": "uint8" },
      { "name": "r", "type": "bytes32" },
      { "name": "s", "type": "bytes32" }
    ],
    "name": "validateOrder_",
    "outputs": [{ "name": "", "type": "bool" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "side", "type": "uint8" },
      { "name": "saleKind", "type": "uint8" },
      { "name": "basePrice", "type": "uint256" },
      { "name": "extra", "type": "uint256" },
      { "name": "listingTime", "type": "uint256" },
      { "name": "expirationTime", "type": "uint256" }
    ],
    "name": "calculateFinalPrice",
    "outputs": [{ "name": "", "type": "uint256" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "protocolFeeRecipient",
    "outputs": [{ "name": "", "type": "address" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [],
    "name": "renounceOwnership",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "addrs", "type": "address[7]" },
      { "name": "uints", "type": "uint256[9]" },
      { "name": "feeMethod", "type": "uint8" },
      { "name": "side", "type": "uint8" },
      { "name": "saleKind", "type": "uint8" },
      { "name": "howToCall", "type": "uint8" },
      { "name": "calldata", "type": "bytes" },
      { "name": "replacementPattern", "type": "bytes" },
      { "name": "staticExtradata", "type": "bytes" }
    ],
    "name": "hashOrder_",
    "outputs": [{ "name": "", "type": "bytes32" }],
    "payable": false,
    "stateMutability": "pure",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "addrs", "type": "address[14]" },
      { "name": "uints", "type": "uint256[18]" },
      { "name": "feeMethodsSidesKindsHowToCalls", "type": "uint8[8]" },
      { "name": "calldataBuy", "type": "bytes" },
      { "name": "calldataSell", "type": "bytes" },
      { "name": "replacementPatternBuy", "type": "bytes" },
      { "name": "replacementPatternSell", "type": "bytes" },
      { "name": "staticExtradataBuy", "type": "bytes" },
      { "name": "staticExtradataSell", "type": "bytes" }
    ],
    "name": "ordersCanMatch_",
    "outputs": [{ "name": "", "type": "bool" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [
      { "name": "addrs", "type": "address[7]" },
      { "name": "uints", "type": "uint256[9]" },
      { "name": "feeMethod", "type": "uint8" },
      { "name": "side", "type": "uint8" },
      { "name": "saleKind", "type": "uint8" },
      { "name": "howToCall", "type": "uint8" },
      { "name": "calldata", "type": "bytes" },
      { "name": "replacementPattern", "type": "bytes" },
      { "name": "staticExtradata", "type": "bytes" },
      { "name": "orderbookInclusionDesired", "type": "bool" }
    ],
    "name": "approveOrder_",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "registry",
    "outputs": [{ "name": "", "type": "address" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "minimumMakerProtocolFee",
    "outputs": [{ "name": "", "type": "uint256" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "addrs", "type": "address[7]" },
      { "name": "uints", "type": "uint256[9]" },
      { "name": "feeMethod", "type": "uint8" },
      { "name": "side", "type": "uint8" },
      { "name": "saleKind", "type": "uint8" },
      { "name": "howToCall", "type": "uint8" },
      { "name": "calldata", "type": "bytes" },
      { "name": "replacementPattern", "type": "bytes" },
      { "name": "staticExtradata", "type": "bytes" }
    ],
    "name": "hashToSign_",
    "outputs": [{ "name": "", "type": "bytes32" }],
    "payable": false,
    "stateMutability": "pure",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [{ "name": "", "type": "bytes32" }],
    "name": "cancelledOrFinalized",
    "outputs": [{ "name": "", "type": "bool" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "owner",
    "outputs": [{ "name": "", "type": "address" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "exchangeToken",
    "outputs": [{ "name": "", "type": "address" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [
      { "name": "addrs", "type": "address[7]" },
      { "name": "uints", "type": "uint256[9]" },
      { "name": "feeMethod", "type": "uint8" },
      { "name": "side", "type": "uint8" },
      { "name": "saleKind", "type": "uint8" },
      { "name": "howToCall", "type": "uint8" },
      { "name": "calldata", "type": "bytes" },
      { "name": "replacementPattern", "type": "bytes" },
      { "name": "staticExtradata", "type": "bytes" },
      { "name": "v", "type": "uint8" },
      { "name": "r", "type": "bytes32" },
      { "name": "s", "type": "bytes32" }
    ],
    "name": "cancelOrder_",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [
      { "name": "addrs", "type": "address[14]" },
      { "name": "uints", "type": "uint256[18]" },
      { "name": "feeMethodsSidesKindsHowToCalls", "type": "uint8[8]" },
      { "name": "calldataBuy", "type": "bytes" },
      { "name": "calldataSell", "type": "bytes" },
      { "name": "replacementPatternBuy", "type": "bytes" },
      { "name": "replacementPatternSell", "type": "bytes" },
      { "name": "staticExtradataBuy", "type": "bytes" },
      { "name": "staticExtradataSell", "type": "bytes" },
      { "name": "vs", "type": "uint8[2]" },
      { "name": "rssMetadata", "type": "bytes32[5]" }
    ],
    "name": "atomicMatch_",
    "outputs": [],
    "payable": true,
    "stateMutability": "payable",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "addrs", "type": "address[7]" },
      { "name": "uints", "type": "uint256[9]" },
      { "name": "feeMethod", "type": "uint8" },
      { "name": "side", "type": "uint8" },
      { "name": "saleKind", "type": "uint8" },
      { "name": "howToCall", "type": "uint8" },
      { "name": "calldata", "type": "bytes" },
      { "name": "replacementPattern", "type": "bytes" },
      { "name": "staticExtradata", "type": "bytes" }
    ],
    "name": "validateOrderParameters_",
    "outputs": [{ "name": "", "type": "bool" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "INVERSE_BASIS_POINT",
    "outputs": [{ "name": "", "type": "uint256" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      { "name": "addrs", "type": "address[14]" },
      { "name": "uints", "type": "uint256[18]" },
      { "name": "feeMethodsSidesKindsHowToCalls", "type": "uint8[8]" },
      { "name": "calldataBuy", "type": "bytes" },
      { "name": "calldataSell", "type": "bytes" },
      { "name": "replacementPatternBuy", "type": "bytes" },
      { "name": "replacementPatternSell", "type": "bytes" },
      { "name": "staticExtradataBuy", "type": "bytes" },
      { "name": "staticExtradataSell", "type": "bytes" }
    ],
    "name": "calculateMatchPrice_",
    "outputs": [{ "name": "", "type": "uint256" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [{ "name": "", "type": "bytes32" }],
    "name": "approvedOrders",
    "outputs": [{ "name": "", "type": "bool" }],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [{ "name": "newOwner", "type": "address" }],
    "name": "transferOwnership",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "name": "registryAddress", "type": "address" },
      { "name": "tokenTransferProxyAddress", "type": "address" },
      { "name": "tokenAddress", "type": "address" },
      { "name": "protocolFeeAddress", "type": "address" }
    ],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "constructor"
  },
  {
    "anonymous": false,
    "inputs": [
      { "indexed": true, "name": "hash", "type": "bytes32" },
      { "indexed": false, "name": "exchange", "type": "address" },
      { "indexed": true, "name": "maker", "type": "address" },
      { "indexed": false, "name": "taker", "type": "address" },
      { "indexed": false, "name": "makerRelayerFee", "type": "uint256" },
      { "indexed": false, "name": "takerRelayerFee", "type": "uint256" },
      { "indexed": false, "name": "makerProtocolFee", "type": "uint256" },
      { "indexed": false, "name": "takerProtocolFee", "type": "uint256" },
      { "indexed": true, "name": "feeRecipient", "type": "address" },
      { "indexed": false, "name": "feeMethod", "type": "uint8" },
      { "indexed": false, "name": "side", "type": "uint8" },
      { "indexed": false, "name": "saleKind", "type": "uint8" },
      { "indexed": false, "name": "target", "type": "address" }
    ],
    "name": "OrderApprovedPartOne",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      { "indexed": true, "name": "hash", "type": "bytes32" },
      { "indexed": false, "name": "howToCall", "type": "uint8" },
      { "indexed": false, "name": "calldata", "type": "bytes" },
      { "indexed": false, "name": "replacementPattern", "type": "bytes" },
      { "indexed": false, "name": "staticTarget", "type": "address" },
      { "indexed": false, "name": "staticExtradata", "type": "bytes" },
      { "indexed": false, "name": "paymentToken", "type": "address" },
      { "indexed": false, "name": "basePrice", "type": "uint256" },
      { "indexed": false, "name": "extra", "type": "uint256" },
      { "indexed": false, "name": "listingTime", "type": "uint256" },
      { "indexed": false, "name": "expirationTime", "type": "uint256" },
      { "indexed": false, "name": "salt", "type": "uint256" },
      { "indexed": false, "name": "orderbookInclusionDesired", "type": "bool" }
    ],
    "name": "OrderApprovedPartTwo",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [{ "indexed": true, "name": "hash", "type": "bytes32" }],
    "name": "OrderCancelled",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      { "indexed": false, "name": "buyHash", "type": "bytes32" },
      { "indexed": false, "name": "sellHash", "type": "bytes32" },
      { "indexed": true, "name": "maker", "type": "address" },
      { "indexed": true, "name": "taker", "type": "address" },
      { "indexed": false, "name": "price", "type": "uint256" },
      { "indexed": true, "name": "metadata", "type": "bytes32" }
    ],
    "name": "OrdersMatched",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [{ "indexed": true, "name": "previousOwner", "type": "address" }],
    "name": "OwnershipRenounced",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      { "indexed": true, "name": "previousOwner", "type": "address" },
      { "indexed": true, "name": "newOwner", "type": "address" }
    ],
    "name": "OwnershipTransferred",
    "type": "event"
  }
]

`
