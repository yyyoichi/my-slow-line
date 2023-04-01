import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

export class TokenizedFetch {
  readonly instance: AxiosInstance;
  constructor() {
    this.instance = axios.create({
      baseURL: '/api',
      withCredentials: true,
    });
  }
  async safe(options: AxiosRequestConfig) {
    const res = await this.instance.get('/safe');
    const token = res.headers['x-csrf-token'];
    options.headers = {
      ...options.headers,
      'X-Csrf-Token': token,
    };
    return await this.instance(options);
  }
}
