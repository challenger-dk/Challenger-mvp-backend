-- Seed initial EULA version 1.0 for en-EN locale
-- This is a data migration, not a schema migration
-- You can edit the content below and update the hash accordingly

INSERT INTO "eula_versions" ("version", "locale", "content", "content_hash", "is_active", "created_at")
VALUES (
  '1.0',
  'en-EN',
  'End User License Agreement (EULA)
Challengers
Last updated: 21/01/2026

This End User License Agreement ("Agreement") is a legally binding agreement between you ("User") and Challengers ("we", "us", or "our") governing your use of the Challengers mobile application ("App").

By downloading, installing, or using the App, you agree to be bound by this Agreement.

1. Age Requirement
The App is intended for users 13 years of age or older.
By using the App, you confirm that you meet this age requirement.

2. License
We grant you a limited, non-exclusive, non-transferable, revocable license to use the App for personal and non-commercial purposes in accordance with this Agreement.

3. User Account
To use the App, you must create a user account. You are responsible for:
- ensuring that the information you provide is accurate
- maintaining the confidentiality of your login credentials
- all activity that occurs under your account

4. User Generated Content (UGC)
The App allows users to create and share content, including but not limited to:
- profile information and profile pictures
- messages and chats
- challenges, teams, invitations, and related content

By uploading or sharing content through the App, you represent and warrant that:
- you own or have the necessary rights to the content
- the content does not infringe the rights of any third party
- the content is not illegal, threatening, hateful, harassing, sexually explicit, or otherwise objectionable

You grant Challengers a non-exclusive, royalty-free license to host, store, display, and distribute your content to the extent necessary to operate and improve the App.

5. Reporting, Blocking, and Moderation
Users may:
- report other users
- report challenges
- block other users, causing their content to be filtered from the user''s experience

Reported content is reviewed manually, and we aim to respond within 24 hours.

We reserve the right to:
- remove content that violates this Agreement
- suspend or terminate user accounts in cases of repeated or serious violations

6. Prohibited Conduct
You agree not to:
- upload unlawful, misleading, or abusive content
- harass, threaten, or abuse other users
- impersonate another person
- attempt to bypass security, authentication, or technical limitations
- use the App for spam or unauthorized commercial purposes

7. Responsibility for Content
All user-generated content is the sole responsibility of the user who created it.
Challengers is not responsible for user content but reserves the right to take action where violations occur.

8. Privacy
Personal data is processed in accordance with our Privacy Policy, which forms an integral part of this Agreement.

9. Intellectual Property
The App, including its code, design, graphics, logos, and functionality, is owned by Challengers and is protected by applicable intellectual property laws.

10. Disclaimer
The App is provided on an "as is" and "as available" basis.
We make no warranties regarding:
- uninterrupted or error-free operation
- the App meeting specific expectations or requirements

11. Limitation of Liability
To the maximum extent permitted by law, Challengers shall not be liable for:
- indirect or consequential damages
- loss of data
- damages arising from user-generated content or interactions between users

12. Termination
We may suspend or terminate your access to the App at any time if you violate this Agreement.

13. Changes to This Agreement
We reserve the right to update this Agreement. Users will be notified of material changes via the App.

14. Governing Law
This Agreement shall be governed by and construed in accordance with Danish law.

15. Contact
If you have any questions regarding this Agreement, please contact us at:
📧 support@challengers-app.com',
  '8685d80d88364428e6a590882af4fd51bb703825feb46a130d9e58e4b4903386',
  true,
  NOW()
)
ON CONFLICT ("version", "locale") DO NOTHING;

-- Seed initial EULA version 1.0 for da-DK locale
INSERT INTO "eula_versions" ("version", "locale", "content", "content_hash", "is_active", "created_at")
VALUES (
  '1.0',
  'da-DK',
  'End User License Agreement (EULA)
Challengers
Sidst opdateret: 21/01/2026

Denne End User License Agreement ("Aftalen") er en juridisk bindende aftale mellem dig ("Brugeren") og Challengers ("vi", "os" eller "vores") vedrørende brugen af mobilapplikationen Challengers ("Appen").

Ved at downloade, installere eller bruge Appen accepterer du denne Aftale.

1. Alderskrav
Appen er beregnet til brugere på 13 år eller derover.
Ved at bruge Appen bekræfter du, at du opfylder dette alderskrav.

2. Licens
Vi giver dig en begrænset, ikke-eksklusiv, ikke-overdragelig og tilbagekaldelig licens til at bruge Appen til personlig og ikke-kommerciel brug i overensstemmelse med denne Aftale.

3. Brugerkonto
For at anvende Appen skal du oprette en brugerkonto. Du er ansvarlig for:
- at de oplysninger, du angiver, er korrekte
- at beskytte dine loginoplysninger
- al aktivitet, der sker via din konto

4. User Generated Content (UGC)
Appen giver mulighed for, at brugere kan oprette og dele indhold, herunder (men ikke begrænset til):
- profiloplysninger og profilbilleder
- beskeder og chats
- challenges, hold, invitationer og relateret indhold

Ved at uploade eller dele indhold i Appen erklærer du, at:
- du har rettighederne til indholdet
- indholdet ikke krænker tredjeparts rettigheder
- indholdet ikke er ulovligt, truende, hadefuldt, chikanerende, seksuelt eksplicit eller på anden måde stødende

Du giver Challengers en ikke-eksklusiv, royalty-fri licens til at hoste, gemme, vise og distribuere dit indhold i det omfang, det er nødvendigt for at drive og forbedre Appen.

5. Rapportering, blokering og moderation
Brugere har mulighed for at:
- rapportere andre brugere
- rapportere challenges
- blokere andre brugere, så deres indhold filtreres fra brugerens oplevelse

Rapporteret indhold bliver gennemgået manuelt, og vi bestræber os på at reagere inden for 24 timer.

Vi forbeholder os retten til:
- at fjerne indhold, der overtræder denne Aftale
- at suspendere eller lukke brugerkonti ved gentagne eller grove overtrædelser

6. Forbudt adfærd
Det er ikke tilladt at:
- uploade ulovligt, misvisende eller krænkende indhold
- chikanere, true eller misbruge andre brugere
- udgive sig for at være en anden person
- forsøge at omgå sikkerhed, authentication eller tekniske begrænsninger
- anvende Appen til spam eller uautoriserede kommercielle formål

7. Ansvar for indhold
Alt brugerindhold er den enkelte brugers ansvar.
Challengers er ikke ansvarlig for indhold, som brugere deler, men forbeholder sig retten til at gribe ind ved overtrædelser.

8. Privatliv
Behandling af personoplysninger sker i overensstemmelse med vores Privatlivspolitik, som er en integreret del af denne Aftale.

9. Immaterielle rettigheder
Appen, herunder kode, design, grafik, logoer og funktionalitet, tilhører Challengers og er beskyttet af gældende ophavsret og immaterielle rettigheder.

10. Ansvarsfraskrivelse
Appen leveres "som den er" og "som tilgængelig".
Vi giver ingen garantier for:
- uafbrudt eller fejlfri drift
- at Appen altid opfylder specifikke forventninger

11. Ansvarsbegrænsning
Challengers kan ikke holdes ansvarlig for:
- indirekte tab
- datatab
- tab som følge af brugerindhold eller interaktioner mellem brugere

12. Opsigelse
Vi kan til enhver tid suspendere eller opsige din adgang til Appen, hvis denne Aftale overtrædes.

13. Ændringer
Vi forbeholder os retten til at opdatere denne Aftale. Ved væsentlige ændringer vil brugere blive informeret via Appen.

14. Gældende lov
Denne Aftale er underlagt dansk ret.

15. Kontakt
Har du spørgsmål til denne Aftale, kan du kontakte os på:
📧 support@challengers-app.com',
  '317f0399e78927d2522b9fdb63409b9114434367e0607064c71c8d1bc3226c1a',
  true,
  NOW()
)
ON CONFLICT ("version", "locale") DO NOTHING;
