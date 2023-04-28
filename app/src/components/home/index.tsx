import Link from 'next/link';
import React from 'react';
import { UiHead } from '../frame';

export type HeadProps = {
  isLogined: boolean;
  account: UserProps;
};
export const Head = (props: HeadProps) => {
  return <UiHead>{props.isLogined ? <User name={props.account.name} /> : <Guest />}</UiHead>;
};

type UserProps = {
  name: string;
};
const User = ({ name }: UserProps) => {
  return (
    <div className='ml-auto mr-4 flex aspect-square w-10 items-center justify-center rounded-full bg-my-white p-1 pb-2 text-xl font-bold text-my-black'>
      {name[0]}
    </div>
  );
};

const Guest = () => {
  return (
    <div className='ml-auto px-4'>
      <Link prefetch={true} className='mx-2 inline-block rounded border-2 border-my-white px-3 py-1' href={'/login'}>
        login
      </Link>
      <Link
        prefetch={true}
        className='mx-2 inline-block rounded border-2 border-my-white bg-my-white px-3 py-1 font-semibold text-my-black'
        href={'/signin'}
      >
        signin
      </Link>
    </div>
  );
};
