import {rateLimit} from 'express-rate-limit';

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 min
  limit: 100,
  standardHeaders: 'draft-7',
  legacyHeaders: false,
});

export {limiter};
