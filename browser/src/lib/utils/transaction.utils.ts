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
	const signatureId = edge.node.tags.find(tag => tag.name === 'Forward-For')?.value;
	const signatureType = edge.node.tags.find(tag => tag.name === 'Signature-Type')?.value;
	if(signatureType){
		const signatureTypeInt = parseInt(signatureType);
		if (signatureTypeInt ==1){
			if (signatureId && signature && eventId){
				try {
					const res = ed25519.verify(signature, eventId, signatureId);
					if (!res){
						console.log('Signature is invalid');
					}
					return res;
				} catch (error) {
					console.log('Error verifying signature:', error);
					return false;
				}
			}else{
				console.log('Missing signatureId, signature, eventId, or userId');
				return false;
			}
		}else if (signatureTypeInt == 2){
			if (signatureId && signature && eventId && userId){
				try {
					const provider = new ethers.JsonRpcProvider('https://polygon-rpc.com') as unknown as Provider;
					const res = verifyMessage({
						signer: signatureId,
						message: eventId,
						signature: signature,
						// this is needed so that smart contract signatures can be verified; this property can also be a viem PublicClient
						provider,
					})
					if (!res){
						console.log('Signature is invalid');
					}
					return res;
				} catch (error) {
					console.log('Error verifying signature:', error);
					return false;
				}
			}else{
				console.log('Missing signatureId, signature, eventId, or userId');
				return false;
			}
		}else if (signatureTypeInt == 3){
			if (eventId && userId && signature){
				try {
					const res = schnorr.verify(signature, eventId, userId);
					if (!res){
						console.log('Signature is invalid');
					}
					return res;
				} catch (error) {
					console.log('Error verifying signature:', error);
					return false;
				}
			}else{
				console.log('Missing signatureId, signature, eventId, or userId');
				return false;
				
			}
		}else{
			console.log('Signature type is unsupported or None');
			return false;
		}
	}
	console.log('Signature type is missing');
	return false;

	
};
