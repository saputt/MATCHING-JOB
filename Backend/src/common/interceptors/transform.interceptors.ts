import {
  CallHandler,
  ExecutionContext,
  HttpStatus,
  Injectable,
  NestInterceptor,
} from '@nestjs/common';
import { map, Observable } from 'rxjs';

export interface Response<T> {
  success: boolean;
  statusCode: number;
  message: string;
  data: T;
}

@Injectable()
export class TransformInterceptor<
  T extends Record<string, any>,
> implements NestInterceptor<T, Response<T>> {
  intercept(
    context: ExecutionContext,
    next: CallHandler<T>,
  ): Observable<Response<T>> {
    const ctx = context.switchToHttp();
    const response = ctx.getResponse();

    const statusCode = response.statusCode;

    let defaultMessage = 'Operation successful';
    if (statusCode == HttpStatus.CREATED)
      defaultMessage = 'Resource created successful';

    return next.handle().pipe(
      map((result) => {
        const customMessage = result?.message || defaultMessage;
        const data = result?.data !== undefined ? result.data : result;

        return {
          success: true,
          statusCode: statusCode,
          message: customMessage,
          data: data === result && result?.message ? null : data,
        };
      }),
    );
  }
}
