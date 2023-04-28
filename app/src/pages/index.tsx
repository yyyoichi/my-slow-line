import dynamic from 'next/dynamic';

export default dynamic(() => import('templates/home/'), { ssr: false });
