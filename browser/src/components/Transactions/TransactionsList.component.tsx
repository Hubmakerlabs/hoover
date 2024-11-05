import { ApiResponse } from '@/lib/types';
import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { TransactionCard } from './TransactionCard.component';
import Typography from '@mui/material/Typography';
import Stack from '@mui/material/Stack';
import Button from '@mui/material/Button';

export const TransactionsList = () => {
	const [cursor, setCursor] = useState<string | null>(null);

	const { data } = useQuery({
		queryKey: ['transactions', cursor],
		queryFn: async () => {
			const queryObject = {
				query: `
					query Transactions {
						transactions(
							first: 25,
							${cursor ? `after: "${cursor}"` : ''}
							sort: HEIGHT_DESC,
							tags: [
								{
									name: "App-Name",
									values: ["com.hubmakerlabs.hoover"]
								}
							]
						)
						{
							pageInfo {
								hasNextPage
							}
							edges {
								cursor
								node {
									id
									signature
									recipient
									owner {
										address
									}
									tags {
										name
										value
									}
								}
							}
						}
					}
				`
			};

			const response = await fetch('https://arweave.net/graphql', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(queryObject)
			});

			return response.json() as Promise<ApiResponse>;
		},
		refetchOnWindowFocus: false
	});

	return (
		<>
			<Stack gap={2} py={4} alignItems="center" justifyContent="flex-start">
				{data?.data.transactions.edges.map(edge => (
					<TransactionCard key={edge.node.id} edge={edge} />
				))}
			</Stack>
			<Stack alignItems="flex-end" bottom={40} right={40} position="fixed" gap={1}>
				<Typography variant="body1">{data?.data.transactions.edges.length} transactions</Typography>
				<Button
					variant="contained"
					disabled={!data?.data.transactions.pageInfo.hasNextPage}
					onClick={() => {
						const nextCursor = data?.data.transactions.edges[data.data.transactions.edges.length - 1].cursor;
						setCursor(nextCursor ?? null);
					}}
				>
					Forwards
				</Button>
			</Stack>
		</>
	);
};
