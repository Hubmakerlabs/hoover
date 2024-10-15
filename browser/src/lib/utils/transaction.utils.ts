import { Edge } from '@/lib/types';
import { schnorr} from '@noble/curves/secp256k1';
import { ed25519 } from '@noble/curves/ed25519';
import * as ethers from 'ethers';
import { verifyMessage } from '@ambire/signature-validator';
import type { Provider } from '@ethersproject/providers';

export const isTxValid = (edge: Edge) => {
	const eventId = edge.node.tags.find(tag => tag.name === 'Event-Id')?.value;
	const userId = edge.node.tags.find(tag => tag.name === 'User-Id')?.value;
	const signature = edge.node.tags.find(tag => tag.name === 'Signature')?.value;
	const signatureId = edge.node.tags.find(tag => tag.name === 'Signature-Id')?.value;
	const signatureType = edge.node.tags.find(tag => tag.name === 'Signature-Type')?.value;
	

	if (signatureType && signatureId && signature && eventId && userId){
		if (parseInt(signatureType) == 1){
			try {
				const res = ed25519.verify(signature, eventId, userId);
				return res;
			} catch (error) {
				console.log('Error verifying signature:', error);
				return false;
			}
		} else if (parseInt(signatureType) == 2){
			
			try {
				const provider = new ethers.JsonRpcProvider('https://polygon-rpc.com') as unknown as Provider;
				const res = verifyMessage({
					signer: signatureId,
					message: eventId,
					signature: signature,
					// this is needed so that smart contract signatures can be verified; this property can also be a viem PublicClient
					provider,
				})
				return res;
			} catch (error) {
				console.log('Error verifying signature:', error);
				return false;
			}
		}else{
			return false;
		}
	} else if (eventId && userId && signature){
		try {
			const res = schnorr.verify(signature, eventId, userId);
			return res;
		} catch (error) {
			console.log('Error verifying signature:', error);
			return false;
		}
	}else{
		return false;
	}

	
};
