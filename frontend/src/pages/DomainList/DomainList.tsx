import Breadcrumb from '../../components/Breadcrumb';
import useAuthToken from '../../features/auth/hooks';
import { useAppDispatch, useAppSelector } from '../../app/hooks';
import { useEffect } from 'react';
import { FetchStatus } from '../../app/types';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { domainFetched } from '../../features/domain/domainSlice';
import { encode } from 'js-base64';
import { Domain } from '../../features/domain/types';
import { fetchServerDetails, selectServerDetails, selectServerDetailsFetchStatus } from '../../features/server/serverSlice';
import moment from 'moment';
import { getSiteCertExpiredDays, isSelfSignedCertificate } from '../../features/certificate/utils';
import { Badge } from 'flowbite-react';
import empty from '../../images/empty.png';
import { CERT_ABOUT_TO_EXPIRE_DAYS } from '../../features/certificate/constants';
import Loader from '../../components/Loader/Loader';
import { getWebServerIcon } from '../../features/server/utils';

const emptyPlaceholder = '----------';

const DomainList = () => {
    const { guid } = useParams();
    const [authToken] = useAuthToken();
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const { pathname } = useLocation();

    const serverDetailsSelectStatus = useAppSelector(selectServerDetailsFetchStatus);
    const serverDetails = useAppSelector(selectServerDetails);

    const isLoading = serverDetailsSelectStatus === FetchStatus.Pending;

    useEffect(() => {
        if (
            authToken
            && serverDetailsSelectStatus !== FetchStatus.Pending
        ) {
            dispatch(fetchServerDetails({
                guid: guid as string,
                token: authToken,
            }));
        }
    }, [authToken]);

    const handleDomainClick = (event: React.MouseEvent, domain: Domain): void => {
        event.preventDefault();
        event.stopPropagation();

        dispatch(domainFetched(domain));
        navigate(`${pathname}/${encode(domain.servername)}`);
    };

    const domains = serverDetails?.domains || [];

    return (
        isLoading
            ? <Loader />
            : <>
                <Breadcrumb pageName='Domains' />
                <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                    <div className="border-b border-stroke px-4 py-5 dark:border-strokedark md:px-6 xl:px-7.5">
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
                                <a href='#' key={domain.servername} onClick={(event: React.MouseEvent) => handleDomainClick(event, domain)}>
                                    <div className="p-3 -mx-3 flex items-center gap-3 hover:bg-[#F8FAFD] dark:hover:bg-meta-4 hover:rounded">
                                        <div className="w-8/12 md:w-5/12 xl:w-4/12 flex gap-4 items-center">
                                            <div className="2xsm:h-11 2xsm:w-full 2xsm:max-w-11 2xsm:rounded-full">
                                                <img src={getWebServerIcon(domain.webserver)} />
                                            </div>
                                            <span className="font-medium text-black dark:text-white">{domain.servername}</span>
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
                                </a>
                            );
                        })}
                    </div>
                </div >
            </>
    );
}

export default DomainList;
