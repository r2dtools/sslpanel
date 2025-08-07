import React from 'react';
import Breadcrumb from '../../components/Breadcrumb';
import { Button } from 'flowbite-react';
import {
    HiCheckCircle,
    HiArchiveBox,
    HiMiniArrowPath,
    HiMiniArrowRight,
} from 'react-icons/hi2';
import Code from '../../components/Help/Code';

const workDir = '/opt/r2dtools'

const UpdateSslBot: React.FC = () => {
    return (
        <div>
            <Breadcrumb pageName='Update SSLBot' />
            <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                <div className='flex flex-col gap-10 p-4 sm:p-6 xl:p-9'>
                    <div>
                        <h3 className="mb-5 text-xl font-bold text-black dark:text-white flex gap-2 items-center">
                            <HiCheckCircle className='text-blue-700' />Step 1: Connect to Your Server
                        </h3>
                        <p className="font-medium mb-2">
                            Use SSH to access your server:
                        </p>
                        <Code>{`ssh your-user@your-server-ip`}</Code>
                    </div>
                    <div>
                        <h3 className="mb-5 text-xl font-bold text-black dark:text-white flex gap-2 items-center">
                            <HiArchiveBox className='text-blue-700' />Step 2: Update the SSLBot
                        </h3>
                        <p className="font-medium mb-2">
                            Download and unpack the latest SSLBot archive:
                        </p>
                        <Code>{
                            `wget -O sslbot.tar.gz https://github.com/r2dtools/sslbot/releases/latest/download/r2dtools-sslbot.tar.gz \\
&& mkdir -p ${workDir} \\
&& tar -xzf sslbot.tar.gz -C ${workDir}`
                        }</Code>
                    </div>
                    <div>
                        <h3 className="mb-5 text-xl font-bold text-black dark:text-white flex gap-2 items-center">
                            <HiMiniArrowPath className='text-blue-700' />Step 3: Restart the SSLBot service
                        </h3>
                        <Code>
                            {`systemctl restart sslbot.service`}
                        </Code>
                        <p className="font-medium mt-2 mb-2">
                            To confirm the agent is running, use:
                        </p>
                        <Code>{`systemctl status sslbot.service`}</Code>
                        <p className="font-medium mt-2">
                            You should see the service status as active
                        </p>
                    </div>
                    <div className='text-center'>
                        <Button href='/servers' color='blue' className='mt-7.5 inline-flex'>Get Started<HiMiniArrowRight className="ml-2 h-5 w-5" /></Button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default UpdateSslBot;
