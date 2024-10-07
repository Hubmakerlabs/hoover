import { useQuery } from '@tanstack/react-query';
import Arweave from 'arweave';
import { Edge, Tag } from '@/lib/types';
import { isTxValid } from '@/lib/utils';
import Typography from '@mui/material/Typography';
import Skeleton from '@mui/material/Skeleton';
import Stack from '@mui/material/Stack';
// import Divider from '@mui/material/Divider';
import Tooltip from '@mui/material/Tooltip';
// import Box from '@mui/material/Box';
import { useTheme } from '@mui/material/styles';

export const TransactionCard = ({ edge }: { edge: Edge }) => {
	const theme = useTheme();

	return (
		<Stack
			gap={2}
			justifyContent="space-between"
			alignItems="flex-start"
			p={2}
			borderRadius={1}
			width="100%"
			sx={{
				boxShadow: 'rgba(255, 255, 255, 0.24) 0px 3px 8px'
			}}
		>
			<Stack direction="row" gap={2} alignItems="center" width="100%" justifyContent="space-between">
				<Typography
					variant="h5"
					sx={{
						textAlign: 'left',
						overflow: 'hidden',
						textOverflow: 'ellipsis',
						whiteSpace: 'nowrap'
					}}
				>
					{/*Event-ID: <b>{edge.node.id}</b>*/}
				</Typography>
				{isTxValid(edge) && (
					<Tooltip title="Verified">
						<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none">
							<path
								d="m8.38 12 2.41 2.42 4.83-4.84"
								stroke={theme.palette.success.main}
								strokeWidth="1.5"
								strokeLinecap="round"
								strokeLinejoin="round"
							/>
							<path
								d="M10.75 2.45c.69-.59 1.82-.59 2.52 0l1.58 1.36c.3.26.86.47 1.26.47h1.7c1.06 0 1.93.87 1.93 1.93v1.7c0 .39.21.96.47 1.26l1.36 1.58c.59.69.59 1.82 0 2.52l-1.36 1.58c-.26.3-.47.86-.47 1.26v1.7c0 1.06-.87 1.93-1.93 1.93h-1.7c-.39 0-.96.21-1.26.47l-1.58 1.36c-.69.59-1.82.59-2.52 0l-1.58-1.36c-.3-.26-.86-.47-1.26-.47H6.18c-1.06 0-1.93-.87-1.93-1.93V16.1c0-.39-.21-.95-.46-1.25l-1.35-1.59c-.58-.69-.58-1.81 0-2.5l1.35-1.59c.25-.3.46-.86.46-1.25V6.2c0-1.06.87-1.93 1.93-1.93h1.73c.39 0 .96-.21 1.26-.47l1.58-1.35Z"
								stroke={theme.palette.success.main}
								strokeWidth="1.5"
								strokeLinecap="round"
								strokeLinejoin="round"
							/>
						</svg>
					</Tooltip>
				)}
			</Stack>
			<Stack alignItems="flex-start" width="100%">
				<Typography
					variant="body1"
					sx={{
						textAlign: 'left',
						overflow: 'hidden',
						textOverflow: 'ellipsis',
						whiteSpace: 'nowrap',
						width: '100%'
					}}
				>
					{/*User-Id: <b>{edge.node.owner.address || '-'}</b>*/}
				</Typography>
				<Typography
					variant="body1"
					sx={{
						textAlign: 'left',
						overflow: 'hidden',
						textOverflow: 'ellipsis',
						whiteSpace: 'nowrap',
						width: '100%'
					}}
				>
				</Typography>
				{/*<Divider flexItem sx={{ my: 1 }} />*/}
				{edge.node.tags
					.filter(tag => tag.name !== 'App-Name' && tag.name !== 'App-Version')
					.map(tag => (
						<TxTag key={tag.name + '-' + tag.value} tag={tag} />
					))}
				{/*<Divider flexItem sx={{ my: 1 }} />*/}
				<TxData id={edge.node.id} />
			</Stack>
		</Stack>
	);
};

const TxTag = ({ tag }: { tag: Tag }) => {
	return (
		<Tooltip title={tag.value}>
			<Typography
				variant="body1"
				sx={{
					textAlign: 'left',
					overflow: 'hidden',
					textOverflow: 'ellipsis',
					whiteSpace: 'nowrap',
					width: '100%'
				}}
			>
				{tag.name}:{' '}
				<b>
					{tag.value} {tag.name === 'Unix-Time' && `(${new Date(Number(tag.value) * 1000).toLocaleString()})`}
				</b>
			</Typography>
		</Tooltip>
	);
};

const TxData = ({ id }: { id: string }) => {
	const { data, isLoading, error } = useQuery({
		queryKey: ['transaction', id],
		queryFn: async () => {
			const arweave = Arweave.init({				host: '127.0.0.1',
				port: 1984,
				protocol: 'http'
			});

			const tx = await arweave.transactions.getData(id, { decode: true });
			const jsonString = new TextDecoder().decode(tx as Uint8Array);

			return JSON.parse(jsonString);
		}
	});

	if (isLoading) {
		return <Skeleton variant="rounded" width="95%" height={200} />;
	}

	if (error) {
		return (
			<Typography variant="body1" color="error.main">
				{/*Error loading TX data: <b>{error.message}</b>*/}
			</Typography>
		);
	}

	return (
		<>
			{data.Content && (
				<Typography
					variant="body1"
					sx={{
						wordBreak: 'break-all'
					}}
				>
					<b>Content:</b>
					<div>{data.Content}</div>
				</Typography>
			)}
			{/*<Box*/}
			{/*	component="pre"*/}
			{/*	sx={{*/}
			{/*		textAlign: 'left',*/}
			{/*		overflowY: 'auto',*/}
			{/*		whiteSpace: 'pre-wrap',*/}
			{/*		width: '95%',*/}
			{/*		maxHeight: 200,*/}
			{/*		padding: 2,*/}
			{/*		borderRadius: 1,*/}
			{/*		backgroundColor: 'primary.main'*/}
			{/*	}}*/}
			{/*>*/}
			{/*	{JSON.stringify(data, null, 2)}*/}
			{/*</Box>*/}
		</>
	);
};
