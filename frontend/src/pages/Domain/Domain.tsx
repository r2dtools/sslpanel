import { Avatar, Badge, Button, ToggleSwitch, Tooltip } from "flowbite-react";
import Breadcrumb from "../../components/Breadcrumb";
import Loader from "../../components/Loader/Loader";
import { HiMiniEye, HiMiniLockClosed } from "react-icons/hi2";
import { useParams } from 'react-router-dom';
import { decode } from 'js-base64';
import Error404 from '../Error404';
import { HiOutlineQuestionMarkCircle } from 'react-icons/hi2';
import { useAppDispatch, useAppSelector } from "../../app/hooks";
import { fetchServerDomain, secureServerDomain, selectDomain, selectDomainFetchStatus, selectDomainSecureStatus } from "../../features/server/domainSlice";
import { FetchStatus } from "../../app/types";
import useAuthToken from "../../features/auth/hooks";
import moment from "moment";
import {
    getCertificateIssuerIcon,
    getCertificateIssuerCode,
    getSiteCertExpiredDays,
    getWebServerIcon,
} from "../../features/server/utils";
import { CERT_ABOUT_TO_EXPIRE_DAYS } from "../../features/server/constants";
import { useEffect, useState } from "react";
import SecureDomainDrawer from "../../features/server/components/SecureDomainDrawer";
import { DomainSecurePayload } from "../../features/server/types";

const emptyPlaceholder = '----------';

const isDnsNameSecure = (certificateDnsNames: string[], name: string): boolean =>
    Boolean(certificateDnsNames.find(certificateDnsName => certificateDnsName === name));

const Domain = () => {
    const { name } = useParams();
    const { guid } = useParams();
    const [authToken] = useAuthToken();
    const dispatch = useAppDispatch();
    const domainName = decode(name || '');
    const domainSelectStatus = useAppSelector(selectDomainFetchStatus);
    const domainSecureStatus = useAppSelector(selectDomainSecureStatus);
    const domain = useAppSelector(selectDomain);
    const [secureFormOpen, setSecureFormOpen] = useState(false);

    if (!domainName) {
        return <Error404 />
    }

    useEffect(() => {
        if (authToken && domainSelectStatus !== FetchStatus.Pending) {
            dispatch(fetchServerDomain({ guid: guid as string, domainName, token: authToken }));
        }
    }, [authToken, domainName, guid]);

    const confPathParts = (domain?.filepath || '').split('/');
    const confFile = confPathParts.pop() || emptyPlaceholder;
    const addresses = (domain?.addresses || []).map(({ host, port }) => `${host}:${port}`);
    const aliases = domain?.aliases || [];
    const certificate = domain?.certificate || null;
    const certificateDnsNames = certificate?.dnsnames || [];
    const certificateExpireDays = getSiteCertExpiredDays(certificate?.validto);
    const expireDuration = certificateExpireDays && certificateExpireDays > 0 ? moment.duration(certificateExpireDays, 'days').humanize() : null;
    const issuerCode = getCertificateIssuerCode(certificate);
    const issuerImg = getCertificateIssuerIcon(issuerCode);
    const organizations = certificate?.organization || [];

    const isLoading = domainSelectStatus === FetchStatus.Pending;

    const handleSubmitSecureForm = async (payload: DomainSecurePayload) => {
        return dispatch(secureServerDomain(payload));
    };

    const handleSecureFormOpen = (): void => {
        setSecureFormOpen(true);
    };

    const handleSecureFormClose = (): void => {
        setSecureFormOpen(false);
    };

    if (!isLoading && !domain) {
        return <Error404 />
    }

    return (
        !isLoading ?
            <>
                <Breadcrumb pageName='Domain' />
                <div className='flex flex-col gap-9'>
                    <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                        <div className="border-b border-stroke px-4 py-5 dark:border-strokedark md:px-6 xl:px-9">
                            <div className="items-center sm:flex">
                                <div className="w-full flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
                                    <Avatar img={getWebServerIcon(domain?.webserver)} size='lg' rounded className='justify-start'>
                                        <div className="space-y-2 font-medium dark:text-white">
                                            <h3 className="inline-block text-2xl font-medium text-black hover:text-primary dark:text-white">{domain?.servername || emptyPlaceholder}</h3>
                                            <div>
                                                {
                                                    domain?.servername && isDnsNameSecure(certificateDnsNames, domain.servername)
                                                        ? <Badge color='success' size='sm' className='inline'>Secure</Badge>
                                                        : <Badge color='failure' size='sm' className='inline'>Insecure</Badge>
                                                }
                                            </div>
                                        </div>
                                    </Avatar>
                                    <div className='flex gap-3'>
                                        <Button color='blue' onClick={handleSecureFormOpen}>
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
                                            {
                                                issuerImg && <div><img src={issuerImg} className='max-w-[200px]' /></div>
                                            }
                                            <div>
                                                <dt>Issuer CN</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    {certificate?.issuer.cn || emptyPlaceholder}
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Issuer Organization</dt>
                                                <dd className="font-bold text-black dark:text-white flex flex-col gap-2">
                                                    {
                                                        organizations.length
                                                            ? organizations.map(organization => (
                                                                <div key={organization}>{organization}</div>
                                                            ))
                                                            : emptyPlaceholder
                                                    }
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Common Name</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    {certificate?.cn || emptyPlaceholder}
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Alternative Names</dt>
                                                <dd className="font-bold text-black dark:text-white flex flex-col gap-2">
                                                    {
                                                        certificateDnsNames.length
                                                            ? certificateDnsNames
                                                                .filter(certificateDnsName => certificateDnsName !== certificate?.cn)
                                                                .map(certificateDnsName => <div key={certificateDnsName}>{certificateDnsName}</div>)
                                                            : emptyPlaceholder
                                                    }
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Expires</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    <span>{certificate?.validto ? moment(certificate.validto).format('LL') : emptyPlaceholder}</span>
                                                    {
                                                        expireDuration !== null && certificateExpireDays !== null
                                                            ? <Badge color={certificateExpireDays < CERT_ABOUT_TO_EXPIRE_DAYS ? 'warning' : 'success'} className='inline ml-2'>{expireDuration}</Badge>
                                                            : null
                                                    }
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
                                                    {
                                                        domain?.ssl
                                                            ? <Badge color='success' className='inline'>Enabled</Badge>
                                                            : <Badge color='failure' className='inline'>Disabled</Badge>
                                                    }
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Aliases</dt>
                                                <dd className="font-bold text-black dark:text-white flex flex-col gap-2">
                                                    {
                                                        aliases.length
                                                            ? aliases.map(
                                                                alias => <div key={alias}>
                                                                    <span>{alias}</span>
                                                                    {
                                                                        isDnsNameSecure(certificateDnsNames, alias)
                                                                            ? <Badge color='success' className='inline ml-2'>Secure</Badge>
                                                                            : <Badge color='failure' className='inline ml-2'>Insecure</Badge>
                                                                    }
                                                                </div>
                                                            )
                                                            : emptyPlaceholder
                                                    }
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Configuration</dt>
                                                <dd className="font-bold text-black dark:text-white flex items-center gap-2">
                                                    <span>{confFile}</span>
                                                    <div className='cursor-pointer'>
                                                        <Tooltip content='show'>
                                                            <HiMiniEye />
                                                        </Tooltip>
                                                    </div>
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>IP Addresses</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    {addresses.length ? addresses.join(', ') : emptyPlaceholder}
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Web Server</dt>
                                                <dd className="font-bold text-black dark:text-white">{domain?.webserver || emptyPlaceholder}</dd>
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
                {
                    domain && authToken && guid && (
                        <SecureDomainDrawer
                            authToken={authToken}
                            guid={guid}
                            domain={domain}
                            open={secureFormOpen}
                            onClose={handleSecureFormClose}
                            onSubmit={handleSubmitSecureForm}
                            loading={domainSecureStatus === FetchStatus.Pending}
                        />)
                }
            </> : <Loader />
    );
}

export default Domain;
