import { Badge, Button, Drawer } from 'flowbite-react';
import { HiMiniKey } from 'react-icons/hi2';
import { Certificate } from '../types';
import { getCertificateIssuerCode, getCertificateIssuerIcon, getSiteCertExpiredDays } from '../utils';
import moment from 'moment';
import { CERT_ABOUT_TO_EXPIRE_DAYS } from '../constants';
import empty from '../../../images/empty.png';

const emptyPlaceholder = '----------';

type CertificateDetailsProps = {
    certificateData: { name: string, certificate: Certificate } | null;
    open: boolean;
    onClose: () => void;
};

const CertificateDetailsDrawer: React.FC<CertificateDetailsProps> = ({
    certificateData,
    open,
    onClose,
}) => {
    const certFileName = certificateData?.name || null;
    const certificate = certificateData?.certificate || null;
    const certificateDnsNames = certificate?.dnsnames || [];
    const certificateAliases = certificateDnsNames.filter(certificateDnsName => certificateDnsName !== certificate?.cn)
    const certificateExpireDays = getSiteCertExpiredDays(certificate?.validto);
    const expireDuration = certificateExpireDays && certificateExpireDays > 0 ? moment.duration(certificateExpireDays, 'days').humanize() : null;
    const issuerCode = getCertificateIssuerCode(certificate);
    const issuerImg = getCertificateIssuerIcon(issuerCode);
    const organizations = certificate?.organization || [];
    const emails = certificate?.emails || [];
    const countries = certificate?.country || [];
    const provinces = certificate?.province || [];
    const localities = certificate?.locality || [];

    const handleFormClose = (): void => {
        onClose();
    };

    return (
        <Drawer className='z-[999] min-w-[400px]' open={open} onClose={handleFormClose} position='right'>
            <Drawer.Header title='Certificate details' titleIcon={HiMiniKey} />
            <Drawer.Items>
                {certificateData ?
                    (
                        <div>
                            <dl className='flex flex-col gap-3'>
                                {
                                    issuerImg && <div><img src={issuerImg} className='max-w-[200px]' /></div>
                                }
                                <div>
                                    <dt>File Name</dt>
                                    <dd className="font-bold text-black dark:text-white">
                                        {`${certFileName}.pem` || emptyPlaceholder}
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
                                            certificateAliases.length
                                                ? certificateAliases.map(certificateDnsName => <div key={certificateDnsName}>{certificateDnsName}</div>)
                                                : emptyPlaceholder
                                        }
                                    </dd>
                                </div>
                                <div>
                                    <dt>Email</dt>
                                    <dd className="font-bold text-black dark:text-white">
                                        {
                                            emails.length
                                                ? emails.map(email => (
                                                    <div key={email}>{email}</div>
                                                ))
                                                : emptyPlaceholder
                                        }
                                    </dd>
                                </div>
                                {
                                    countries.length ? (
                                        <div>
                                            <dt>Country</dt>
                                            <dd className="font-bold text-black dark:text-white">
                                                {
                                                    countries.map(country => (
                                                        <div key={country}>{country}</div>
                                                    ))
                                                }
                                            </dd>
                                        </div>) : null
                                }
                                {
                                    provinces.length ? (
                                        <div>
                                            <dt>Province</dt>
                                            <dd className="font-bold text-black dark:text-white">
                                                {
                                                    provinces.map(province => (
                                                        <div key={province}>{province}</div>
                                                    ))
                                                }
                                            </dd>
                                        </div>) : null
                                }
                                {
                                    localities.length ? (
                                        <div>
                                            <dt>Locality</dt>
                                            <dd className="font-bold text-black dark:text-white">
                                                {
                                                    localities.map(locality => (
                                                        <div key={locality}>{locality}</div>
                                                    ))
                                                }
                                            </dd>
                                        </div>) : null
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
                                <div>
                                    <dd>
                                        <Button color='blue' onClick={handleFormClose}>Close</Button>
                                    </dd>
                                </div>
                            </dl>
                        </div>
                    ) : (
                        <div className='m-auto max-w-24 mt-3'><img src={empty} /></div>
                    )
                }
            </Drawer.Items>
        </Drawer >
    );
};

export default CertificateDetailsDrawer;
