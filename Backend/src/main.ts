import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { BadRequestException, ValidationPipe } from '@nestjs/common';
import { TransformInterceptor } from './common/interceptors/transform.interceptors';
import { GlobalExceptionFilter } from './common/filters/http-exception.filter';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  app.setGlobalPrefix("api")

  app.useGlobalPipes(
    new ValidationPipe({
      whitelist : true,
      transform : true,
      forbidNonWhitelisted : true,
      exceptionFactory: (validationErrors) => {
        const messages = validationErrors.map(
          (error) =>
            `${error.property}: ${Object.values(error.constraints ?? {}).join(', ')}`,
        );
        return new BadRequestException(messages);
      },
    })
  )

  app.useGlobalInterceptors(new TransformInterceptor())

  app.useGlobalFilters(new GlobalExceptionFilter())

  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
