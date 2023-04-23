import { AxiosRequestConfig } from 'axios';
import { tokenizeFetch } from './axiosInstance';

/**@return {number} userId */
export const postSignin = async (email: string, password: string, name = '') => {
  const config: AxiosRequestConfig = {
    data: {
      email,
      password,
      name,
    },
    method: 'POST',
    url: 'signin',
  };
  try {
    const data = await tokenizeFetch<number>(config);
    return data;
  } catch (e) {
    throw e;
  }
};

/**@param id userID */
export const postVerificateCode = async (id: number, code: string) => {
  const config: AxiosRequestConfig = {
    data: {
      id,
      code,
    },
    method: 'POST',
    url: 'codein',
  };
  try {
    await tokenizeFetch<unknown>(config);
    return true;
  } catch (e) {
    throw e;
  }
};
