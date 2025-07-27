import { Popover, Spinner, ToggleSwitch } from 'flowbite-react';
import { FC } from 'react';
import { HiOutlineQuestionMarkCircle } from 'react-icons/hi2';
import { ServerSettings as ServerSettingsType } from '../types';

type ServerSettingsProps = {
    settings: ServerSettingsType;
    onCertbotStatusChange: (value: boolean) => void;
    certbotStatusLoading: boolean;
};

const ServerSettings: FC<ServerSettingsProps> = ({ settings, certbotStatusLoading, onCertbotStatusChange }) => {
    const handleCertbotStatusChange = (value: boolean) => {
        onCertbotStatusChange(value);
    };

    return (
        <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
            <div className="p-4 md:p-6 xl:p-7.5">
                <div className="flex items-start justify-between">
                    <div>
                        <h2 className="text-xl font-bold text-black dark:text-white">Settings</h2>
                    </div>
                </div>
            </div>
            <div className='border-b border-stroke px-4 dark:border-strokedark md:px-6 xl:px-7.5' />
            <div className='p-4 md:p-6 xl:p-7.5'>
                <div className='flex max-w-md flex-col items-start gap-4'>
                    <div className='flex items-center gap-2'>
                        <ToggleSwitch
                            color='green'
                            checked={settings.certbotStatus}
                            label='Enable Certbot Integration'
                            onChange={handleCertbotStatusChange}
                        />
                        <Popover placement='right' trigger='click' content={
                            <div className="text-sm text-gray-500 dark:text-gray-400">
                                <div className="px-3 py-2">
                                    <div className="w-64 p-3">
                                        <p>Allows you to use <a className='font-bold text-blue-600' href='https://certbot.eff.org/' target='_blank'>Certbot</a> to issue and install SSL/TLS certificates on your server. Please make sure <span className='font-bold'>Certbot</span> is already installed before using this feature.</p>
                                    </div>
                                </div>
                            </div>
                        }>
                            <div className="cursor-pointer">
                                <HiOutlineQuestionMarkCircle />
                            </div>
                        </Popover>
                        {certbotStatusLoading ? <Spinner size="sm" /> : null}
                    </div>
                </div>
            </div>
        </div >
    );
};

export default ServerSettings;
