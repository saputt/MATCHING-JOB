import {
  ArgumentsHost,
  Catch,
  ExceptionFilter,
  HttpException,
  HttpStatus,
  Logger,
} from '@nestjs/common';
import { Request, Response } from 'express';

@Catch()
export class GlobalExceptionFilter implements ExceptionFilter {
  private readonly logger = new Logger(GlobalExceptionFilter.name);

  catch(exception: Error | HttpException, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();

    const status =
      exception instanceof HttpException
        ? exception.getStatus()
        : HttpStatus.INTERNAL_SERVER_ERROR;

    let message = 'Internal server error';
    let errors = null;

    if (status === 429) {
      message = 'Too many requests';
    } else if (exception instanceof HttpException) {
      const resResponse = exception.getResponse();
      if (typeof resResponse === 'object' && resResponse !== null) {
        message = (resResponse as any).message || exception.message;
        errors = (resResponse as any).errors || null;
      } else {
        message = exception.message;
      }
    } else if (exception instanceof Error) {
      this.logger.error(
        `Path: ${request.url} | Error: ${exception.message}`,
        exception.stack,
      );
      message = 'Terjadi kesalahan pada sistem internal';
    }

    response.status(status).json({
      success: false,
      statusCode: status,
      path: request.url,
      message: Array.isArray(message) ? message[0] : message,
      errors: errors,
      timestamp: new Date().toISOString(),
    });
  }
}
