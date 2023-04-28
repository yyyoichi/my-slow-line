import dynamic from 'next/dynamic';

export default dynamic(() => import('templates/login/'), { ssr: false });
