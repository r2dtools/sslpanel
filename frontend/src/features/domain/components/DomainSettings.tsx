import { Popover, Spinner, ToggleSwitch } from "flowbite-react";
import { Domain, DomainSettings as DomainSettingsType } from "../types";
import { FC, useState } from "react";
import { HiOutlineQuestionMarkCircle } from "react-icons/hi2";

type DomainSettingsProps = {
    domain: Domain;
    settings: DomainSettingsType;
    commonDirStatusLoading: boolean;
    onCommonDirStatusChange: (value: boolean) => void;
};

const DomainSettings: FC<DomainSettingsProps> = ({ settings, domain, commonDirStatusLoading, onCommonDirStatusChange }) => {
    const [renewalStatus, setRenewalStatus] = useState<boolean>(settings.renewal);
    const [renewalStatusChangeLoading, setRenewalChangeLoading] = useState<boolean>(false);

    const handleCommonDirStatusChanged = async (value: boolean) => {
        await onCommonDirStatusChange(value);
    };

    const handleRenewalStatusChanged = (value: boolean) => {
        setRenewalStatus(value);
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
            <div className="border-b border-stroke px-4 dark:border-strokedark md:px-6 xl:px-7.5"></div>
            <div className="p-4 md:p-6 xl:p-7.5">
                <div className="flex max-w-md flex-col items-start gap-4">
                    {
                        commonDirStatusLoading
                            ? <Spinner />
                            : (
                                <div className='flex items-center gap-2'>
                                    <ToggleSwitch
                                        color='green'
                                        checked={settings.commondirstatus.status}
                                        label="Enable Common ACME Challenge Directory" onChange={handleCommonDirStatusChanged}
                                    />
                                    <Popover placement='right' trigger='click' content={
                                        <div className="text-sm text-gray-500 dark:text-gray-400">
                                            <div className="px-3 py-2">
                                                <div className="w-64 p-3">
                                                    <p><span className="font-bold">Common Challenge Directory</span> significantly decreases a number of cases when Let’s Encrypt SSL/TLS certificates cannot be issued because of incompatible domain configurations.</p>
                                                </div>
                                            </div>
                                        </div>
                                    }>
                                        <div className="cursor-pointer">
                                            <HiOutlineQuestionMarkCircle className='inline' />
                                        </div>
                                    </Popover>
                                </div>
                            )
                    }
                    <div>
                        {renewalStatusChangeLoading
                            ? <Spinner />
                            : (
                                <ToggleSwitch
                                    color='green'
                                    checked={renewalStatus}
                                    label="Automatic Renewal"
                                    onChange={handleRenewalStatusChanged}
                                />
                            )
                        }
                    </div>

                </div>
            </div>
        </div >
    );
};

export default DomainSettings;
