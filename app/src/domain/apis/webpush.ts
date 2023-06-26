import { AxiosRequestConfig } from 'axios';
import { normalFetch, tokenizeFetch } from './axiosInstance';

export const getVapidPublicKey = async () => {
  const config: AxiosRequestConfig = {
    method: 'GET',
    url: 'me/webpush/vapid_public_key',
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

export const setSubscription = async (subscription: PushSubscriptionJSON) => {
  const config: AxiosRequestConfig = {
    method: 'POST',
    url: 'me/webpush/subscription',
    data: {
      endpoint: subscription.endpoint || '',
      p256hd: subscription.keys?.p256dh || '',
      auth: subscription.keys?.auth || '',
      expirationTime: subscription.expirationTime,
      userAgent: navigator.userAgent,
    },
  };
  try {
    const res = await tokenizeFetch<null>(config);
    return res.status === 200;
  } catch (e) {
    return Error('failture');
  }
};
