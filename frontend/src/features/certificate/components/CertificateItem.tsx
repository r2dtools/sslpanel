import { Link } from 'react-router-dom';
import { Badge, Spinner, Tooltip } from 'flowbite-react';
import { HiMiniCloudArrowDown } from 'react-icons/hi2';
import {
    getCertificateIssuerCode,
    getCertificateIssuerIcon,
    getSiteCertExpiredDays,
    isSelfSignedCertificate,
} from '../utils';
import sslIcon from '../../../images/certificate/ca.png';
import moment from 'moment';
import { CERT_ABOUT_TO_EXPIRE_DAYS } from '../constants';
import { Certificate } from '../types';
import { useState } from 'react';

interface CertificateItemProps {
    name: string
    storage: string
    certificate: Certificate
    onCertificateDownload: (name: string, storage: string) => Promise<any>
    onClick: (certificate: Certificate, name: string) => void
};

const CertificateItem: React.FC<CertificateItemProps> = ({ name, storage, certificate, onCertificateDownload, onClick }) => {
    const [actionLoading, setActionLoading] = useState<boolean>(false);
    const fileName = `${name}.pem`

    const handleDownload = async (event: React.MouseEvent<SVGElement, MouseEvent>) => {
        setActionLoading(true);

        event.preventDefault();
        event.stopPropagation();

        await onCertificateDownload(name, storage);

        setActionLoading(false);
    };

    const handleClick = (event: React.MouseEvent<HTMLAnchorElement>): void => {
        preventClick(event);
        onClick(certificate, name)
    };

    const preventClick = (event: React.MouseEvent) => {
        event.preventDefault();
        event.stopPropagation();
    };

    const code = getCertificateIssuerCode(certificate);
    const icon = getCertificateIssuerIcon(code);

    const validTo = certificate.validto;
    const expireDays = getSiteCertExpiredDays(validTo);
    const expireDuration = expireDays && expireDays > 0 ? moment.duration(expireDays, 'days').humanize() : null;
    const dnsnames = [...new Set([certificate.cn].concat(certificate.dnsnames || []))];

    const isSelfSigned = isSelfSignedCertificate(certificate);

    return (
        <Link to="#" onClick={handleClick}>
            <div className="p-3 flex items-center gap-3 hover:bg-[#F8FAFD] dark:hover:bg-meta-4 hover:rounded">
                <div className="w-3/12 md:w-4/12">
                    {
                        icon ? (
                            <div className="max-w-40">
                                <img src={icon} className='w-full' />
                            </div>
                        ) : (
                            <div className='flex items-center gap-4'>
                                <div className="2xsm:h-11 2xsm:max-w-11">
                                    <img src={sslIcon} />
                                </div>
                                <div className='flex flex-col gap-1'>
                                    <span className="font-bold">Default CA</span>
                                    {isSelfSigned && <Badge color='warning' className='inline'>Self-Signed</Badge>}
                                </div>
                            </div>
                        )
                    }
                </div>
                <div className="w-4/12 flex flex-col gap-1 font-medium">
                    {dnsnames.map((name: string) => <div className='truncate text-black dark:text-white' key={name}>{name}</div>)}
                    <div className='text-sm truncate'>{`[${fileName}, ${storage} storage]`}</div>
                </div>
                <div className="w-3/12 md:flex md:flex-col gap-1 lg:gap-2 md:items-center lg:flex-row">
                    <span className="font-medium hidden md:block text-black dark:text-white">
                        {moment(validTo).format('LL')}
                    </span>
                    {expireDays && expireDuration ? (
                        <Badge className='inline' color={expireDays < CERT_ABOUT_TO_EXPIRE_DAYS ? 'warning' : 'success'}>
                            {expireDuration}
                        </Badge>
                    ) : null}
                    {expireDays !== null && expireDays <= 0 ? (
                        <Badge className='inline' color='failure'>Expired</Badge>
                    ) : null}
                </div>
                <div className="w-2/12 text-center lg:w-1/12" onClick={preventClick}>
                    {actionLoading
                        ? <Spinner />
                        : (
                            <button className="flex justify-between mx-auto gap-2">
                                <Tooltip content="Download" >
                                    <HiMiniCloudArrowDown size={20} className='hover:text-red-500' onClick={handleDownload} />
                                </Tooltip>
                            </button>
                        )
                    }
                </div>
            </div>
        </Link>
    );
};

export default CertificateItem;
