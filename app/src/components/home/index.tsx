import Link from 'next/link';
import React from 'react';
import { UiHead } from '../frame';
import { DivProps } from 'components/atoms/type';
import { NonNullablePick } from 'components';

export type HeadProps = {
  isLogined: boolean;
  user: UserProps;
};
export const Head = (props: HeadProps) => {
  return <UiHead>{props.isLogined ? <User {...props.user} /> : <Guest />}</UiHead>;
};

type UserProps = {
  icon: NonNullablePick<DivProps, 'children'>;
  logout: NonNullablePick<DivProps, 'onClick'>;
};
const User = ({ icon, logout }: UserProps) => {
  const logoutProps: DivProps = {
    ...logout,
    className: 'mx-2 inline-block rounded border-2 border-my-white px-3 py-1',
    children: 'logout',
  };
  const iconProps: DivProps = {
    ...icon,
    className:
      'flex aspect-square w-10 items-center justify-center rounded-full bg-my-white p-1 pb-2 text-xl font-bold text-my-black',
  };
  return (
    <div className='ml-auto flex gap-4'>
      <div {...logoutProps} />
      <div {...iconProps} />
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
