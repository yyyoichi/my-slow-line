import { AxiosRequestConfig } from 'axios';
import { normalFetch } from './axiosInstance';

export const getVapidPublicKey = async () => {
  const config: AxiosRequestConfig = {
    method: 'GET',
    url: 'webpush/vapid_public_key',
  };
  let key: string;
  try {
    const res = await normalFetch<string>(config);
    key = res.status === 200 ? res.data : '';
  } catch (e) {
    return Error('failture');
  }
  return key;
};
