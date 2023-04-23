import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

/**fetch api in csrf-safe */
export const tokenizeFetch = async <T>(config: AxiosRequestConfig) => {
  const instance = axios.create({
    baseURL: '/api',
  });
  const res = await instance.get('/safe');
  const token = res.headers['x-csrf-token'];
  config.headers = {
    ...config.headers,
    'X-Csrf-Token': token,
  };
  const response = await instance.request<T>(config);
  const data = await response.data;
  if (response.status !== axios.HttpStatusCode.Ok) throw new Error(data as string);
  return await response.data;
};
