import { Edge } from '@/lib/types';
import { schnorr} from '@noble/curves/secp256k1';
import { ed25519 } from '@noble/curves/ed25519';
import * as ethers from 'ethers';
import { verifyMessage } from '@ambire/signature-validator';
import type { Provider } from '@ethersproject/providers';
import { IdResolver } from '@atproto/identity';
import { verifySignature } from '@atproto/crypto'
import { CID } from 'multiformats/cid';

export const isTxValid = (edge: Edge) => {
	const eventId = edge.node.tags.find(tag => tag.name === 'Event-Id')?.value;
	const userId = edge.node.tags.find(tag => tag.name === 'User-Id')?.value;
	const signature = edge.node.tags.find(tag => tag.name === 'Signature')?.value;
	const signatureId = edge.node.tags.find(tag => tag.name === 'Signer')?.value;
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
		} else if (signatureTypeInt == 4) {
			if (userId && signature && eventId) {
				console.log('Signature type is Bluesky');
				const idResolver = new IdResolver();
		
				const verifyBlueskySignature = async (): Promise<boolean> => {
					try {
						// Parse the CID
						const cid = CID.parse(eventId);
						const hashBytes = cid.multihash.digest; // Extract the multihash
		
						// Convert the signature from hex to Uint8Array
						if (signature.length % 2 !== 0) {
							throw new Error('Hex string must have an even number of characters');
						}
						const byteArray = new Uint8Array(signature.length / 2);
						for (let i = 0; i < signature.length; i += 2) {
							byteArray[i / 2] = parseInt(signature.slice(i, i + 2), 16);
						}
		
						// Resolve the public key using the DID
						const publicKey = await idResolver.did.resolveAtprotoKey(userId, true);
		
						// Verify the signature
						const isValid = verifySignature(publicKey, hashBytes, byteArray); // Custom function
						if (!isValid) {
							console.log('Signature is invalid');
							return false;
						}
		
						console.log('Signature is valid');
						return true;
					} catch (error) {
						console.log('Error verifying Bluesky signature:', error);
						return false;
					}
				};
		
				// Initiate signature verification
				return verifyBlueskySignature()
					.then((isValid) => isValid)
					.catch((error) => {
						console.log('Unhandled error during Bluesky signature verification:', error);
						return false;
					});
			} else {
				console.log('Missing signatureId, signature, or eventId');
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
