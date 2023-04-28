import { AxiosRequestConfig } from 'axios';
import { normalFetch } from './axiosInstance';

type Me = {
  id: string;
  name: string;
  email: string;
  loginAt: string;
  createAt: string;
  updateAt: string;
};

export const getMe = async () => {
  const config: AxiosRequestConfig = {
    method: 'GET',
    url: '/users/me',
  };
  let me: MyAccount;
  try {
    const res = await normalFetch<Me>(config);
    me = res.status === 200 ? new MyAccount(res.data) : new MyAccount();
  } catch (e) {
    me = new MyAccount();
  }
  return me;
};

export class MyAccount {
  readonly has: boolean;
  readonly id: string;
  readonly name: string;
  readonly email: string;
  readonly loginAt: Date;
  readonly createAt: Date;
  readonly updateAt: Date;
  constructor(me?: Me) {
    this.has = Boolean(me);
    this.id = me?.id || '';
    this.name = me?.name || '';
    this.email = me?.email || '';
    this.loginAt = new Date(me?.loginAt || '');
    this.createAt = new Date(me?.createAt || '');
    this.updateAt = new Date(me?.updateAt || '');
  }
}
