import { Badge } from 'flowbite-react';
import { Link, useLocation } from 'react-router-dom';
import { Domain } from '../types';
import React from 'react';
import empty from '../../../images/empty.png';
import moment from 'moment';
import { encode } from 'js-base64';
import { getSiteCertExpiredDays, isSelfSignedCertificate } from '../utils';
import { CERT_ABOUT_TO_EXPIRE_DAYS } from '../constants';

type ServerDomainListProps = {
    domains: Domain[];
};

const emptyPlaceholder = '----------';

const ServerDomainList: React.FC<ServerDomainListProps> = ({ domains }) => {
    const { pathname } = useLocation();

    return (
        <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
            <div className="p-4 md:p-6 xl:p-7.5">
                <div className="flex items-start justify-between">
                    <div>
                        <h2 className="text-xl font-bold text-black dark:text-white">Domains</h2>
                    </div>
                </div>
            </div>
            <div className="border-b border-stroke px-4 pb-5 dark:border-strokedark md:px-6 xl:px-7.5">
                <div className="flex items-center gap-3">
                    <div className="w-8/12 md:w-5/12 xl:w-4/12">
                        <span className="font-medium">Domain</span>
                    </div>
                    <div className="w-4/12 md:w-3/12 xl:w-2/12">
                        <span className="font-medium">Status</span>
                    </div>
                    <div className="hidden xl:block w-3/12">
                        <span className="font-medium">Certificate</span>
                    </div>
                    <div className="hidden w-3/12 md:w-4/12 xl:w-3/12 md:block">
                        <span className="font-medium">Expires</span>
                    </div>
                </div>
            </div>
            <div className="p-4 md:p-6 xl:p-7.5">
                <div className="mx-auto max-w-[80px]">
                    {!domains.length && <img src={empty} />}
                </div>
                {domains.map((domain: Domain) => {
                    const cert = domain.certificate;
                    const issuer = cert?.issuer;
                    const validTo = cert?.validto;
                    const organization = (issuer?.organization || []).join(', ');
                    const expireDays = getSiteCertExpiredDays(validTo || null);
                    const expireDuration = expireDays && expireDays > 0 ? moment.duration(expireDays, 'days').humanize() : null;

                    return (
                        <Link to={`${pathname}/domain/${encode(domain.servername)}`} key={domain.servername}>
                            <div className="p-3 -mx-3 flex items-center gap-3 hover:bg-[#F8FAFD] dark:hover:bg-meta-4 hover:rounded">
                                <div className="w-8/12 md:w-5/12 xl:w-4/12">
                                    <span className="font-medium">{domain.servername}</span>
                                </div>
                                <div className="w-4/12 md:w-3/12 xl:w-2/12">
                                    {cert ? (
                                        <Badge color={cert?.isvalid ? 'success' : 'failure'} size='sm' className='inline'>
                                            {cert?.isvalid ? 'Valid' : 'Invalid'}
                                        </Badge>) : emptyPlaceholder
                                    }
                                </div>
                                <div className="hidden xl:block w-3/12">
                                    <span className="font-medium">{organization || issuer?.cn || emptyPlaceholder}</span>
                                    {isSelfSignedCertificate(cert) && (
                                        <Badge color='warning' className='inline ml-2'>Self-Signed</Badge>
                                    )}
                                </div>
                                <div className="hidden w-3/12 md:w-4/12 xl:w-3/12 md:block">
                                    <span className="font-medium">
                                        {validTo ? moment(validTo).format('LL') : emptyPlaceholder}
                                    </span>
                                    {expireDays && expireDuration ? (
                                        <Badge className='inline ml-2' color={expireDays < CERT_ABOUT_TO_EXPIRE_DAYS ? 'warning' : 'success'}>
                                            {expireDuration}
                                        </Badge>
                                    ) : null}
                                    {expireDays !== null && expireDays <= 0 ? (
                                        <Badge className='inline ml-2' color='failure'>Expired</Badge>
                                    ) : null}
                                </div>
                            </div>
                        </Link>
                    );
                })}
            </div>
        </div >
    );
};

export default ServerDomainList;
