import { AxiosRequestConfig } from 'axios';
import { APIERROR, tokenizeFetch } from './axiosInstance';

export const putParticipantAt = async (body: PostSessionKeyBody) => {
  const config: AxiosRequestConfig<PostSessionKeyBody> = {
    method: 'POST',
    url: `me/sessionkey`,
    data: body,
  };
  try {
    const res = await tokenizeFetch<null>(config);
    if (res.status === 200) {
      return true;
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

type PostSessionKeyBody = {
  sessionID: number;
  inviteeID: number;
  key: string;
};
