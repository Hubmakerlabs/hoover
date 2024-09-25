import { Edge } from '@/lib/types';
import { schnorr } from '@noble/curves/secp256k1';

export const isTxValid = (edge: Edge) => {
	const eventId = edge.node.tags.find(tag => tag.name === 'Event-Id')?.value;
	const userId = edge.node.tags.find(tag => tag.name === 'User-Id')?.value;
	const signature = edge.node.tags.find(tag => tag.name === 'Signature')?.value;

	if (!eventId || !userId || !signature) {
		return false;
	}

	try {
		const res = schnorr.verify(signature, eventId, userId);
		return res;
	} catch (error) {
		console.log('Error verifying signature:', error);
		return false;
	}
};
