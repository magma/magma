import {rateLimit} from 'express-rate-limit';
import {RATE_LIMIT_CONFIG} from '../../../config/config';

const limiter = rateLimit({...RATE_LIMIT_CONFIG});

export {limiter};
