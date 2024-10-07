export interface Tag {
	name: string;
	value: string;
}

export interface Node {
	id: string;
	signature: string;
	recipient: string;
	owner: {
		address: string;
	};
	tags: Tag[];
}

export interface Edge {
	cursor: string;
	node: Node;
}

export interface PageInfo {
	hasNextPage: boolean;
}

export interface Transactions {
	pageInfo: PageInfo;
	edges: Edge[];
}

export interface Data {
	transactions: Transactions;
}

export interface ApiResponse {
	data: Data;
}
