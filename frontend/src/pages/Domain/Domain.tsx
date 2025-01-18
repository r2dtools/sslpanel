import { Avatar, Badge, Button, ToggleSwitch, Tooltip } from "flowbite-react";
import Breadcrumb from "../../components/Breadcrumb";
import Loader from "../../components/Loader/Loader";
import { HiMiniEye, HiMiniLockClosed } from "react-icons/hi2";
import { useParams } from 'react-router-dom';
import { decode } from 'js-base64';
import Error404 from '../Error404';
import nginxIcon from '../../images/nginx.svg';
import comodo from '../../images/ca/sectigo.svg';
import { HiOutlineQuestionMarkCircle } from 'react-icons/hi2';


const Domain = () => {
    const { name } = useParams();
    const domainName = decode(name || '');

    if (!domainName) {
        return <Error404 />
    }

    const isLoading = false;

    return (
        !isLoading ?
            <>
                <Breadcrumb pageName='Domain' />
                <div className='flex flex-col gap-9'>
                    <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                        <div className="border-b border-stroke px-4 py-5 dark:border-strokedark md:px-6 xl:px-9">
                            <div className="items-center sm:flex">
                                <div className="w-full flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
                                    <Avatar img={nginxIcon} size='lg' rounded className='justify-start'>
                                        <div className="space-y-2 font-medium dark:text-white">
                                            <h3 className="inline-block text-2xl font-medium text-black hover:text-primary dark:text-white">{domainName}</h3>
                                            <div>
                                                <Badge color='success' size='sm' className='inline'>Secure</Badge>
                                            </div>
                                        </div>
                                    </Avatar>
                                    <div className='flex gap-3'>
                                        <Button color='blue'>
                                            <HiMiniLockClosed className="mr-2 h-5 w-4" />
                                            Secure
                                        </Button>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div className="p-4 md:p-6 xl:p-9">
                            <div className='flex flex-col gap-3'>
                                <div className='flex gap-6 flex-col sm:flex-row sm:gap-3'>
                                    <div className='sm:basis-1/2'>
                                        <h4 className="mb-3 font-bold text-black dark:text-white uppercase">Certificate Information</h4>
                                        <dl className='flex flex-col gap-3'>
                                            <div>
                                                <img src={comodo} className='max-w-[200px]' />
                                            </div>
                                            <div>
                                                <dt>Issuer CN</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    R3
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Issuer Organization</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    Let`s Encrypt
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Common Name</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    r2dtools.work.gd
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Alternative Names</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    www.r2dtools.work.gd, webmail.r2dtools.work.gd
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Expires</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    <span>25 Jan 2025</span>
                                                    <Badge color='warning' className='inline ml-2'>30 days</Badge>
                                                </dd>
                                            </div>
                                        </dl>
                                    </div>
                                    <div className='sm:basis-1/2'>
                                        <h4 className="mb-3 font-bold text-black dark:text-white uppercase">Domain Information</h4>
                                        <dl className='flex flex-col gap-3'>
                                            <div>
                                                <dt>SSL</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    <Badge color='success' className='inline'>Enabled</Badge>
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Aliases</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    <span>www.domain.com</span>
                                                    <Badge color='failure' className='inline ml-2'>Insecure</Badge>
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Configuration</dt>
                                                <dd className="font-bold text-black dark:text-white flex items-center gap-2">
                                                    <span>example.com.conf</span>
                                                    <div className='cursor-pointer'>
                                                        <Tooltip content='show'>
                                                            <HiMiniEye />
                                                        </Tooltip>
                                                    </div>
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>IP Addresses</dt>
                                                <dd className="font-bold text-black dark:text-white">*</dd>
                                            </div>
                                            <div>
                                                <dt>Web Server</dt>
                                                <dd className="font-bold text-black dark:text-white">Nginx</dd>
                                            </div>
                                        </dl>
                                    </div>

                                </div>
                            </div>
                        </div>
                    </div>
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
                                <div className='flex items-center gap-2'>
                                    <ToggleSwitch color='green' checked={false} label="Enable Common Challenge Directory" onChange={() => ''} />
                                    <HiOutlineQuestionMarkCircle className='inline' />
                                </div>
                                <div>
                                    <ToggleSwitch color='green' checked label="Automatic Renewal" onChange={() => ''} />
                                </div>

                            </div>
                        </div>
                    </div >
                </div >
            </> : <Loader />
    );
}

export default Domain;
