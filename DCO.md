# DCO

You must sign-off all commits on the originating branch for a PR, which certifies that you wrote it or otherwise have the right to pass it on as an open-source contribution.
The rules are pretty simple: if you can certify the below (from [developercertificate.org](https://developercertificate.org/)):

> Developer's Certificate of Origin 1.1
>
> By making a contribution to this project, I certify that:
>
> (a) The contribution was created in whole or in part by me and I
> have the right to submit it under the open source license
> indicated in the file; or
>
> (b) The contribution is based upon previous work that, to the best
> of my knowledge, is covered under an appropriate open source
> license and I have the right under that license to submit that
> work with modifications, whether created in whole or in part
> by me, under the same open source license (unless I am
> permitted to submit under a different license), as indicated
> in the file; or
>
> (c) The contribution was provided directly to me by some other
> person who certified (a), (b) or (c) and I have not modified
> it.
>
> (d) I understand and agree that this project and the contribution
> are public and that a record of the contribution (including all
> personal information I submit with it, including my sign-off) is
> maintained indefinitely and may be redistributed consistent with
> this project or the open source license(s) involved.

Then you just add a line to every git commit message:

> Signed-off-by: Joe Smith <joe.smith@email.com>

You need to use your real name to contribute (sorry, no pseudonyms or anonymous contributions).
If you set your user.name and user.email git configs, you can sign your commit automatically with [git commit -s](https://git-scm.com/docs/git-commit#Documentation/git-commit.txt--s).
