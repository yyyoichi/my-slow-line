import { AxiosRequestConfig } from 'axios';
import { APIERROR, tokenizeFetch } from './axiosInstance';

export const putParticipantAt = async (sessionID: number, body: PutParticipantAtBody) => {
  const config: AxiosRequestConfig<PutParticipantAtBody> = {
    method: 'POST',
    url: `me/participants/${sessionID}`,
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

type PutParticipantAtBody = {
  userID: number;
  status: ParticipantStatus;
};
