import React from 'react';
import Breadcrumb from '../../components/Breadcrumb';
import { Button, List, ListItem } from 'flowbite-react';
import {
    HiCheckCircle,
    HiArrowDownOnSquare,
    HiArchiveBox,
    HiMiniKey,
    HiMiniArrowPath,
    HiGlobeAlt,
    HiMiniExclamationCircle,
    HiMiniArrowRight,
} from 'react-icons/hi2';
import Code from '../../components/Help/Code';

const InstallSslBot: React.FC = () => {
    return (
        <div>
            <Breadcrumb pageName='Install SSLBot' />
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
                            <HiArrowDownOnSquare className='text-blue-700' />Step 2: Download the SSLBot Installer
                        </h3>
                        <p className="font-medium mb-2">
                            Download the latest version of the SSLBot installer using one of the commands below:
                        </p>
                        <Code>{`wget -O installer https://github.com/r2dtools/sslbot-installer/releases/latest/download/installer`}</Code>
                        <p className="font-medium mb-2 mt-2">
                            Or using curl:
                        </p>
                        <Code>{`curl -o installer https://github.com/r2dtools/sslbot-installer/releases/latest/download/sslbot-installer`}</Code>
                    </div>
                    <div>
                        <h3 className="mb-5 text-xl font-bold text-black dark:text-white flex gap-2 items-center">
                            <HiArchiveBox className='text-blue-700' />Step 3: Install the SSLBot
                        </h3>
                        <p className="font-medium mb-2">
                            Run the installer:
                        </p>
                        <Code>{`chmod +x installer | ./installer install`}</Code>
                        <p className="font-medium mt-2 mb-2">
                            The agent will be installed to the <span className='font-bold'>/opt/sslbot</span> directory
                        </p>
                        <p className="font-medium mb-2">
                            To confirm the agent is running, use:
                        </p>
                        <Code>{`systemctl status sslbot.service`}</Code>
                        <p className="font-medium mt-2">
                            You should see the service status as active
                        </p>
                    </div>
                    <div>
                        <h3 className="mb-5 text-xl font-bold text-black dark:text-white flex gap-2 items-center">
                            <HiMiniKey className='text-blue-700' />Step 4: Generate Token
                        </h3>
                        <Code>{`/opt/sslbot generate-token`}</Code>
                        <p className="font-medium mt-2 mb-2">
                            The token will be saved in the file:
                        </p>
                        <Code>{`/opt/sslbot/config/params.yaml`}</Code>
                    </div>
                    <div>
                        <h3 className="mb-5 text-xl font-bold text-black dark:text-white flex gap-2 items-center">
                            <HiMiniArrowPath className='text-blue-700' />Step 5: Restart the SSLBot
                        </h3>
                        <Code>{`systemctl restart sslbot.service`}</Code>
                    </div>
                    <div>
                        <h3 className="mb-5 text-xl font-bold text-black dark:text-white flex gap-2 items-center">
                            <HiGlobeAlt className='text-blue-700' />Step 6: Register the Server in SSLPanel
                        </h3>
                        <List ordered>
                            <ListItem>Log in to your SSLPanel dashboard</ListItem>
                            <ListItem>Go to the "Servers" page and click the "Add" button</ListItem>
                            <ListItem>
                                Enter the following details:
                                <List nested>
                                    <ListItem><span className='font-bold'>Server Name</span>: (any name you prefer)</ListItem>
                                    <ListItem><span className='font-bold'>IP Address</span>: your server's IP</ListItem>
                                    <ListItem><span className='font-bold'>Port</span>: 60150 (default SSLBot port)</ListItem>
                                    <ListItem><span className='font-bold'>Token</span>: the token generated earlier</ListItem>
                                </List>
                            </ListItem>
                            <ListItem>Click "Save", then open the created server entry</ListItem>
                        </List>
                    </div>
                    <div>
                        <h3 className="mb-5 text-xl font-bold text-black dark:text-white flex gap-2 items-center">
                            <HiMiniExclamationCircle className='text-blue-700' />Final Check
                        </h3>
                        <p className="font-medium mb-2">
                            If everything is correct, the server status will show "active"
                        </p>
                        <p className="font-medium mb-2">
                            If you see a connection error:
                        </p>
                        <List>
                            <ListItem>Double-check the IP address, port, and token</ListItem>
                            <ListItem>Make sure SSLBot is running</ListItem>
                            <ListItem>Ensure no firewall is blocking SSLBot port</ListItem>
                        </List>
                    </div>
                    <div className='text-center'>
                        <Button href='/servers' color='blue' className='mt-7.5 inline-flex'>Get Started<HiMiniArrowRight className="ml-2 h-5 w-5" /></Button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default InstallSslBot;
